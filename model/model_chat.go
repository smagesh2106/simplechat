package model

import (
	"github.com/gorilla/websocket"
)

type ChatRoom struct {
	RoomID string `json:"roomid"`
}

type Message struct {
	room string
	data []byte
}

type Connection struct {
	ws   *websocket.Conn
	send chan []byte
}
type Subscription struct {
	room string
}
