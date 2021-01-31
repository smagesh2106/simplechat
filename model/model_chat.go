package model

import (
	"log"
	//"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	//writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	//pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	//pingPeriod = (pongWait * 9) / 10
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
	User string
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

func (s *Subscription) ReadThread() {
	c := s.Conn
	defer func() {
		Hub.Unregister <- *s
		c.Ws.Close()
	}()

	c.Ws.SetReadLimit(maxMessageSize)
	//c.Ws.SetReadDeadline(time.Now().Add(pongWait))
	//c.Ws.SetPongHandler(func(string) error { c.Ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.Ws.ReadMessage()
		log.Println("Read this msg from client :", msg)
		if err != nil {
			log.Printf("Could not read the msg from ws :%v\n", err)
			break
		}
		Hub.Broadcast <- Message{Room: s.Room, User: s.User, Data: msg}
	}
}

func (s *Subscription) WriteThread() {
	c := s.Conn
	defer func() {
		Hub.Unregister <- *s
		c.Ws.Close()
	}()
	//ticker := time.NewTicker(pingPeriod)

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
			/*
				case <-ticker.C:
					if err := c.Ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
						return
					}
			*/
		}

	}

	/*
		for {
			select {
			case message, ok := <-c.Send:
				if !ok {
					c.write(websocket.CloseMessage, []byte{})
					return
				}
				if err := c.write(websocket.TextMessage, message); err != nil {
					return
				}
			case <-ticker.C:
				if err := c.write(websocket.PingMessage, []byte{}); err != nil {
					return
				}

			}
		}
	*/
}

// write writes a message with the given message type and payload.
/*
func (c *Connection) write(mt int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(mt, payload)
}
*/
