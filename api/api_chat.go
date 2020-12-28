package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	c "github.com/securechat/chat"
	//"path"
	//"path/filepath"
)

func ChatLaunch(w http.ResponseWriter, r *http.Request) {
	//Filename := path.Base(r.URL.String())
	log.Println("Attempting to serve", "index.html")
	//http.ServeFile(w, r, filepath.Join(".", "html", "index.html"))
	http.ServeFile(w, r, "./html/index.html")
}

func ChatSession(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomId := params["roomId"]
	userId := params["userId"]

	if len(userId) > 0 {
		c.ServeWs(w, r, roomId, userId)
	} else {
		c.ServeWs(w, r, roomId, "")
	}
}
