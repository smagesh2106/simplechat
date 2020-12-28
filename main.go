package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	api "github.com/securechat/api"
	com "github.com/securechat/driver"
	mod "github.com/securechat/model"

	"github.com/gorilla/handlers"
	"gopkg.in/validator.v2"
	//"github.com/gorilla/handlers"
)

func main() {
	log.Println("Start Secure Chat ...")

	label := false

	//Initialize mongodb and start, wait till it becomes available.
	for {
		if label {
			break
		}
		// Initializa mongo db.
		err := com.Init_Mongo()
		if err != nil {
			log.Printf("Error setting up mongoDB :%v", err)
			time.Sleep(5 * time.Second)
		} else {
			label = true
		}
	}
	//Init Rooms
	mod.InitRooms()

	serverPort := os.Getenv("SERVER_PORT")
	httpOnly := os.Getenv("HTTP_ONLY")
	httpEnabled, _ := strconv.ParseBool(httpOnly)
	log.Printf("HTTP only :%v", httpOnly)

	router := api.NewRouter()

	// Serve static files from html dir.
	router.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("./html"))))

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	//AllowCredentials: true,
	//Debug:            true,

	//Set up custom validations
	validator.SetValidationFunc("nonemail", mod.Nonemail)
	validator.SetValidationFunc("nonint", mod.Nonint)

	//Set up periodic cleanup.
	go func() {
		for now := range time.Tick(time.Hour * 24) {
			fmt.Println("Cleanup Start :%v", now)
			cleanup()
		}
	}()
	go mod.Hub.Run()
	log.Printf("Server started...")
	if httpEnabled {
		log.Printf("Running in HTTP mode")
		log.Fatal(http.ListenAndServe(":"+serverPort, handlers.CORS(origins, headers, methods)(router)))
		//log.Fatal(http.ListenAndServe(":"+serverPort, handlers.CORS(header, methods, origins)(router)))
	} else {
		serverSSLPort := os.Getenv("SERVER_SSL_PORT")
		serverCrtFile := os.Getenv("SERVER_SSL_CERT")
		serverKeyFile := os.Getenv("SERVER_SSL_KEY")

		if _, err := os.Stat(serverKeyFile); os.IsNotExist(err) {
			log.Panic(err)
		}

		if _, err := os.Stat(serverCrtFile); os.IsNotExist(err) {
			log.Panic(err)
		}

		log.Printf("Running in HTTPS mode")
		log.Fatal(http.ListenAndServeTLS(":"+serverSSLPort, serverCrtFile, serverKeyFile, handlers.CORS(origins, headers, methods)(router)))
	}
}

func cleanup() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in cleanup :", r)
		}
	}()
	//do all clean up logic
	log.Println("Start cleanup.")
	time.Sleep(5 * time.Second)
	log.Println("End cleanup.")
}
