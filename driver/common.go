package driver

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection *mongo.Collection

/*
var Chat_Session *mongo.Collection
var Chat_Summary *mongo.Collection
var ProjectCollection *mongo.Collection
var Categories *mongo.Collection
var UserRegistration *mongo.Collection
var UserSession *mongo.Collection
*/
var Log *log.Logger

func Init_Connections(database string) {
	//<FIXME: Need different collection for diff object types >
	UserCollection = Client.Database(database).Collection("users")
	/*
		Chat_Session = Client.Database(database).Collection("chat_session")
		Chat_Summary = Client.Database(database).Collection("chat_summary")
		ProjectCollection = Client.Database(database).Collection("projects")
		Categories = Client.Database(database).Collection("skill_sub_category")
		UserRegistration = Client.Database(database).Collection("user_registration")
		UserSession = Client.Database(database).Collection("user_session")
	*/
	Log.Println("Connections to tables done..")
}
