package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	mod "github.com/securechat/model"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(w http.ResponseWriter, r *http.Request, room string, user string) {
	log.Println("Inside broadcast router....")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS Socket upgrade error :%v", err)
		return
	}
	c := mod.Connection{Send: make(chan []byte, 256), Ws: ws}
	s := mod.Subscription{Conn: &c, Room: room, User: user}
	log.Printf("Registering client %v", s)
	mod.Hub.Register <- s
	go s.ReadThread()
	go s.WriteThread()
}
