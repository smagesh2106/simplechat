package model

import (
	"context"
	"log"
	"time"

	db "github.com/securechat/driver"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/x/bsonx"
)

type ChatRoom struct {
	RoomID string `validate:"min=3" json:"roomid  bson:"roomid"`
}

type ChatRoomList struct {
	Rooms []ChatRoom `json:"rooms" bson: "rooms"`
}

type ErrMsg struct {
	Error string `json:"error"`
}

/*
 * Chat Model : CreateRoom
 */
func (room *ChatRoom) CreateChatRoom() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mod := mongo.IndexModel{
		//Keys:    bsonx.Doc{{Key: "roomid", Value: bsonx.String("text")}},
		Keys:    bson.D{{"roomid", "text"}},
		Options: options.Index().SetUnique(true),
	}

	_, err := db.ChatRoomDB.Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Printf("Unable to create unique index : %v", err)
		return err
	}
	_, err = db.ChatRoomDB.InsertOne(ctx, room)

	if err != nil {
		log.Printf("Unable to create chat room : %v", err)
		return err
	}
	return nil
}

/*
 * Chat Model : DeleteRoom
 */
func (room *ChatRoom) DeleteChatRoom() (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := db.ChatRoomDB.DeleteOne(ctx, room)

	if err != nil {
		log.Printf("Unable to delete chat room : %v", err)
		return nil, err
	}
	return result, nil
}

/*
 * Chat Model : GetRoomList
 */
func GetAllChatRooms() (*ChatRoomList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.ChatRoomDB.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	if err != nil {
		log.Println("Mongo DB Conn error")
		return nil, err
	}
	roomlist := []ChatRoom{}
	for cursor.Next(ctx) {
		c := ChatRoom{}
		err := cursor.Decode(&c)
		if err != nil {
			return nil, err
		}

		roomlist = append(roomlist, c)
	}

	list := &ChatRoomList{}
	list.Rooms = roomlist
	return list, nil
}

/*
 * Init Rooms
 */
func InitRooms() {
	var list []string = []string{"Pune", "Chennai", "Mumbai", "Bangalore", "Hyderabad", "Punjab", "Delhi", "Ahmedabad", "Kolkata", "Madurai",
		"Belgaum"}

	for _, r := range list {
		room := ChatRoom{RoomID: r}
		(&room).CreateChatRoom()
	}
}
