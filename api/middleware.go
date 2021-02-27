package api

import (
	"encoding/json"
	//"log"
	"net/http"
	//"os"
	//"time"
	//"regexp"
	//"strings"
	u "github.com/securechat/driver"
)

//var log = log.New(os.Stdout, "secure-chat :", log.LstdFlags)

//var regex = regexp.MustCompile("\\s*")

func Validator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not Authorized."))
			return
		} else {
			u.Log.Printf("Token :%v", auth)
			//logger.Println("Token :" + strings.Split(strings.Trim(auth),
		}
		next.ServeHTTP(w, r)
	})
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()
		//logger := log.New(os.Stdout, "iam-policy-administration :", log.LstdFlags)
		u.Log.Printf(
			"%s %s %s", // %s",
			r.Method,
			r.RequestURI,
			name,
		//	time.Since(start),
		)
		inner.ServeHTTP(w, r)
	})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				u.Log.Printf("Generic Error caught :%v", err)

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
