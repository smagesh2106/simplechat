package model

import (
	"log"

	"github.com/gorilla/websocket"
)

const (
	maxMessageSize = 512
)

// Message to the room or user
type Message struct {
	Room       string
	TargetUser string
	Data       []byte
	Owner      string
}

// Websocket connection, read data from websocket and send it to the byte chan.
type Connection struct {
	Ws    *websocket.Conn
	Send  chan []byte
	Owner string
}

// Chat subcription.
type Subscription struct {
	Room       string
	TargetUser string
	Conn       *Connection
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

		case m := <-H.Unicast:
			connections := H.Rooms[m.Room]
			for c := range connections {
				if c.Owner == m.TargetUser || c.Owner == m.Owner {
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
}

func (s *Subscription) ReadThread() {
	c := s.Conn
	defer func() {
		Hub.Unregister <- *s
		c.Ws.Close()
	}()

	c.Ws.SetReadLimit(maxMessageSize)
	for {
		_, msg, err := c.Ws.ReadMessage()
		msg = []byte("Host-User :" + string(msg))
		if err != nil {
			log.Printf("Could not read the msg from ws :%v\n", err)
			break
		}

		if s.TargetUser == "boradcast" {
			Hub.Broadcast <- Message{Room: s.Room, TargetUser: s.TargetUser, Owner: s.Conn.Owner, Data: msg}
		} else {
			Hub.Unicast <- Message{Room: s.Room, TargetUser: s.TargetUser, Owner: s.Conn.Owner, Data: msg}
		}
	}
}

func (s *Subscription) WriteThread() {
	c := s.Conn
	defer func() {
		Hub.Unregister <- *s
		c.Ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			//		c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				if err := c.Ws.WriteMessage(websocket.TextMessage, message); err != nil {
					return
				}
			}
		}

	}
}
