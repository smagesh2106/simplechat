package api

import (
	"net/http"

	//"path/filepath"

	"github.com/gorilla/mux"
	c "github.com/securechat/chat"
	u "github.com/securechat/driver"
	//u "github.com/securechat/util"
	//"path"
	//"path/filepath"
)

//var log = log.New(os.Stdout, "secure-chat :", log.LstdFlags)

func ChatLaunch(w http.ResponseWriter, r *http.Request) {
	//Filename := path.Base(r.URL.String())
	u.Log.Println("Attempting to serve", "index.html")
	//http.ServeFile(w, r, filepath.Join(".", "html", "index.html"))
	http.ServeFile(w, r, "./html/index.html")
	//http.ServeFile(w, r, "index.html")
	u.Log.Println("Done serving ", "index.html")
}

func ChatBroadcast(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomId := params["roomId"]

	u.Log.Println("roomId :" + roomId)
	c.ServeWs(w, r, roomId, "")
}

func ChatUnicast(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomId := params["roomId"]
	userId := params["userId"]

	u.Log.Println("roomId :" + roomId + ",  userId :" + userId)
	c.ServeWs(w, r, roomId, userId)
}
