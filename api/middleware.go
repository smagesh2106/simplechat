package api

import (
	"log"
	"net/http"
	"os"
	//"time"
	//"regexp"
	//"strings"
)

var logger = log.New(os.Stdout, "secure-chat :", log.LstdFlags)

//var regex = regexp.MustCompile("\\s*")

func Validator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println("middleware -1 ")
		auth := r.Header.Get("Authorization")

		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not Authorized."))
			return
		} else {
			logger.Printf("Token :%v", auth)
			//logger.Println("Token :" + strings.Split(strings.Trim(auth),
		}
		next.ServeHTTP(w, r)
	})
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()
		//logger := log.New(os.Stdout, "iam-policy-administration :", log.LstdFlags)
		logger.Println("middleware -2 ")
		logger.Printf(
			"%s %s %s", // %s",
			r.Method,
			r.RequestURI,
			name,
		//	time.Since(start),
		)
		inner.ServeHTTP(w, r)
	})
}
