package driver

import (
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollectionDB *mongo.Collection
var ChatRoomDB *mongo.Collection

/*
var Chat_Session *mongo.Collection
var Chat_Summary *mongo.Collection
var ProjectCollection *mongo.Collection
var Categories *mongo.Collection
var UserRegistration *mongo.Collection
var UserSession *mongo.Collection
*/
//var Log *log.Logger
var Log = log.New(os.Stdout, "secure-chat :", log.LstdFlags)

func Init_Connections(database string) {
	UserCollectionDB = Client.Database(database).Collection("users")
	ChatRoomDB = Client.Database(database).Collection("rooms")

	/*
		Chat_Session = Client.Database(database).Collection("chat_session")
		Chat_Summary = Client.Database(database).Collection("chat_summary")
		ProjectCollection = Client.Database(database).Collection("projects")
		Categories = Client.Database(database).Collection("skill_sub_category")
		UserRegistration = Client.Database(database).Collection("user_registration")
		UserSession = Client.Database(database).Collection("user_session")
	*/
	log.Println("Connections to tables done..")
}
