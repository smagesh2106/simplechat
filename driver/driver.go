package driver

import (
	"bytes"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

/*
 * Init function to initialize mongo DB connection.
 */
func Init_Mongo() error {

	Log = log.New(os.Stdout, "cresent-web :", log.LstdFlags)
	Log.Println("mongodb set start")

	tlsConfig := new(tls.Config)
	tlsConfig.InsecureSkipVerify = false

	key := os.Getenv("DB_SSL_CLIENT_KEY")
	cert := os.Getenv("DB_SSL_CLIENT_CERT")
	cacert := os.Getenv("SSL_CA_CERT")

	sub, err := AddClientCertFromSeparateFiles(tlsConfig, key, cacert, cert, "")

	//sub, err := mgo.AddClientCertFromConcatenatedFile(tlsConfig, keyCert, "")
	if err != nil {
		Log.Printf("Error while reading certificates, %v", err)
		return err
	} else {
		Log.Printf("certificates added...%v\n", sub)
	}
	Log.Printf("Verify Certificate : %s\n", strconv.FormatBool(!tlsConfig.InsecureSkipVerify))

	//<FIXME include replicaset in the production cluster>
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoAddr := strings.Join([]string{mongoHost, mongoPort}, ":")
	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	database := os.Getenv("DB_NAME")
	mongoAuthSource := os.Getenv("MONGO_AUTHSOURCE")
	//replicaSet := "rs0"
	mongourl := "mongodb://" + mongoAddr + "/" + database + "?ssl=true"

	clientOptions := options.Client().ApplyURI(mongourl).
		SetAuth(options.Credential{
			AuthSource: mongoAuthSource, Username: mongoUsername, Password: mongoPassword,
		}).SetTLSConfig(tlsConfig) //.SetReplicaSet("rs0")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	Client, _ = mongo.Connect(ctx, clientOptions)
	err = Client.Ping(ctx, nil)
	if err != nil {
		Log.Printf("mongo connection error %v", err)
		return err
	}
	DB = Client.Database(database)

	Init_Connections(database)
	Log.Println("mongodb set end")
	return nil
}

//ConnectMongoWithoutSSL
func ConnectMongoWithoutSSL() error {

	var (
		mongoURL = "mongodb://localhost:27017"
		err      error
	)
	Log = log.New(os.Stdout, "cresent-web :", log.LstdFlags)
	Log.Println("mongodb set start")

	// Initialize a new mongo client with options
	Client, err = mongo.NewClient(options.Client().ApplyURI(mongoURL))

	// Connect the mongo client to the MongoDB server
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Client.Connect(ctx)

	// Ping MongoDB
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if err = Client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println("could not ping to mongo db service: %v\n", err)
		return err
	}
	DB = Client.Database(os.Getenv("DB_NAME"))

	Init_Connections(os.Getenv("DB_NAME"))
	return nil

}

/*
 * Read Key / Cert from a concatenated file.
 */
func AddClientCertFromConcatenatedFile(cfg *tls.Config, certKeyFile, keyPassword string) (string, error) {
	data, err := ioutil.ReadFile(certKeyFile)
	if err != nil {
		return "", err
	}

	return addClientCertFromBytes(cfg, data, keyPassword)
}

/*
 * Read Key / Cert from seperate files
 */
func AddClientCertFromSeparateFiles(cfg *tls.Config, keyFile, cacert, certFile, keyPassword string) (string, error) {

	err := addCACertFromFile(cfg, cacert)
	if err != nil {
		Log.Printf("Error while reading CA cert from file :%v", err)
		return "", err
	}

	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		Log.Printf("Error whiel reading key file: %v", err)
		return "", err
	}
	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		Log.Printf("Error while reading cert data: %v", err)
		return "", err
	}

	data := append(keyData, '\n')
	data = append(data, certData...)
	return addClientCertFromBytes(cfg, data, keyPassword)
}

/*
 * Read CA Cert file
 */
func addCACertFromFile(cfg *tls.Config, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		Log.Printf("Error while reading CA cert from file: %v", err)
		return err
	}

	certBytes, err := loadCACert(data)
	if err != nil {
		Log.Printf("Error while loading CA cert: %v", err)
		return err
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		Log.Printf("Error parsing certificate: %v", err)
		return err
	}

	if cfg.RootCAs == nil {
		cfg.RootCAs = x509.NewCertPool()
	}

	cfg.RootCAs.AddCert(cert)

	return nil
}

func loadCACert(data []byte) ([]byte, error) {
	var certBlock *pem.Block

	for certBlock == nil {
		if data == nil || len(data) == 0 {
			return nil, errors.New("no CERTIFICATE section found")
		}

		block, rest := pem.Decode(data)
		if block == nil {
			return nil, errors.New("invalid .pem file")
		}

		switch block.Type {
		case "CERTIFICATE":
			certBlock = block
		}
		data = rest
	}
	return certBlock.Bytes, nil
}

func addClientCertFromBytes(cfg *tls.Config, data []byte, keyPasswd string) (string, error) {
	var currentBlock *pem.Block
	var certBlock, certDecodedBlock, keyBlock []byte

	remaining := data
	start := 0
	for {
		currentBlock, remaining = pem.Decode(remaining)
		if currentBlock == nil {
			break
		}

		if currentBlock.Type == "CERTIFICATE" {
			certBlock = data[start : len(data)-len(remaining)]
			certDecodedBlock = currentBlock.Bytes
			start += len(certBlock)
		} else if strings.HasSuffix(currentBlock.Type, "PRIVATE KEY") {
			//<FIXME> if the key file has a password, replace "" with password
			if keyPasswd != "" && x509.IsEncryptedPEMBlock(currentBlock) {
				var encoded bytes.Buffer
				buf, err := x509.DecryptPEMBlock(currentBlock, []byte(keyPasswd))
				if err != nil {
					Log.Printf("Error while decrypting a key using password: %v", err)
					return "", err
				}

				pem.Encode(&encoded, &pem.Block{Type: currentBlock.Type, Bytes: buf})
				keyBlock = encoded.Bytes()
				start = len(data) - len(remaining)
			} else {
				keyBlock = data[start : len(data)-len(remaining)]
				start += len(keyBlock)
			}
		}
	}
	if len(certBlock) == 0 {
		return "", fmt.Errorf("failed to find CERTIFICATE")
	}
	if len(keyBlock) == 0 {
		return "", fmt.Errorf("failed to find PRIVATE KEY")
	}

	cert, err := tls.X509KeyPair(certBlock, keyBlock)
	if err != nil {
		Log.Printf("Error while creating a certificate: %v", err)
		return "", err
	}

	cfg.Certificates = append(cfg.Certificates, cert)
	crt, err := x509.ParseCertificate(certDecodedBlock)
	if err != nil {
		Log.Printf("Error Parsing Certificate: %v", err)
		return "", err
	}
	return crt.Subject.String(), nil
}
