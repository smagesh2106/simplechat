package model

import (
	"context"
	"log"
	"time"

	db "github.com/securechat/driver"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.in/validator.v2"
)

var Session map[string]time.Time

type User struct {
	FirstName string    `json:"firstName" bson:"firstName" validate:"min=3,max=15,nonnil,nonzero"`
	LastName  string    `json:"lastName" bson:"lastName" validate:"min=3,max=15,nonnil,nonzero"`
	UserID    string    `json:"userid" bson:"userid" validate:"min=3,max=30,nonnil,nonzero,regexp=^[A-Za-z0-9!@#-_.]+$"`
	Password  string    `json:password bson:"password" validate:"min=8,max=30,nonnil,nonzero"`
	Email     string    `json:"email" bson:"email" validate:"min=3,max=30,nonnil,nonzero"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	LastLogin time.Time `json:"lastLogin" bson:"lastLogin"`
}

func (u *User) validate() error {
	if err := validator.Validate(u); err != nil {
		log.Println("Error validating user obj :", u)
		return err
	}
	return nil
}

/*
 * User Model : DeleteUser
 */
func (user *User) DeleteUser() (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := db.UserCollectionDB.DeleteOne(ctx, bson.M{"userid": user.UserID})

	if err != nil {
		log.Printf("Unable to delete user : %v", err)
		return nil, err
	}
	return result, nil
}

/*
 * User Model : CreateUser
 */
func (user *User) CreateUser() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"userid", 1}},
			Options: options.Index().SetUnique(true),
		}, {
			Keys:    bson.D{{"email", 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err = db.UserCollectionDB.Indexes().CreateMany(ctx, indexes, opts)
	if err != nil {
		log.Printf("Unable to create index, error : %v\n", err)
		return err
	}
	_, err = db.UserCollectionDB.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Unable to create user : %v", err)
		return err
	}
	return nil
}
