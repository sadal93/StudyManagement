package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

var userCollection = client.Database("test").Collection("users")

type User struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
	TimeLastAssigned int64
	TimeToSend int64
}

func createUserDocument(doc User) {
	insertResult, err := userCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(insertResult)
}


