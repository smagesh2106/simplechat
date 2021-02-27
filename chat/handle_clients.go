package chat

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	u "github.com/securechat/driver"
	mod "github.com/securechat/model"
)

var i int = 1
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(w http.ResponseWriter, r *http.Request, room string, target string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		u.Log.Printf("WS Socket upgrade error :%v", err)
		return
	}
	//Get Client/Owner from session userinfo
	owner := "cli" + strconv.Itoa(i)
	i++
	u.Log.Println("Owner :" + owner)
	c := mod.Connection{Send: make(chan []byte, 256), Ws: ws, Owner: owner}
	var s = &mod.Subscription{Conn: &c, Room: room}
	//
	if len(target) > 0 {
		s.TargetUser = target
	} else {
		s.TargetUser = "boradcast"
	}
	mod.Hub.Register <- *s
	go s.ReadThread()
	go s.WriteThread()
}
