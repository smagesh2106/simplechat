package model

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Message to the room or user
type Message struct {
	Room string
	User string
	Data []byte
}

// Websocket connection, read data from websocket and send it to the byte chan.
type Connection struct {
	Ws   *websocket.Conn
	Send chan []byte
}

// Chat subcription.
type Subscription struct {
	Room string
	Conn *Connection
}

type chatHub struct {
	Broadcast  chan Message
	Unicast    chan Message
	Register   chan Subscription
	Unregister chan Subscription
	Rooms      map[string]map[*Connection]bool
}

var Hub = chatHub{
	Broadcast:  make(chan Message),
	Unicast:    make(chan Message), //for unicast messages
	Register:   make(chan Subscription),
	Unregister: make(chan Subscription),
	Rooms:      make(map[string]map[*Connection]bool),
}

func (H *chatHub) Run() {
	for {
		select {

		case s := <-H.Register:
			//Get connection map from Hub for a given room
			connections := H.Rooms[s.Room]
			if connections == nil {
				//Create connections map if it does not exist
				connections = make(map[*Connection]bool)
				H.Rooms[s.Room] = connections
			}
			//Set the session connection to 'true'
			H.Rooms[s.Room][s.Conn] = true

		case s := <-H.Unregister:
			//Get connection map from Hub for a given room
			connections := H.Rooms[s.Room]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					//remove the subscriber from the Hub's connection map
					delete(connections, s.Conn)
					//close the subcriber's byte channel
					close(s.Conn.Send)
					//if there are no connections in a room then delete the room.
					if len(connections) == 0 {
						delete(H.Rooms, s.Room)
					}

				}
			}

		case m := <-H.Broadcast:
			connections := H.Rooms[m.Room]
			for c := range connections {
				select {
				case c.Send <- m.Data:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(H.Rooms, m.Room)
					}
				}

			}
		}
	}
}
