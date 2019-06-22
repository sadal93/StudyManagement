package main

import (
	pb "Thesis-demo/api"
	"context"
	"log"
)

var userCollection = client.Database("test").Collection("users")

type User struct {
	UserID string
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

func (s *server)  CreateUser(ctx context.Context, user *pb.UserMetaData) (*pb.UserMetaData, error) {

	userDoc:= User{user.UserID, 0, 3600000}
	createUserDocument(userDoc)

	log.Printf("User Created: %v", user.UserID)
	return &pb.UserMetaData{UserID: user.UserID, TimeLastAssigned: user.TimeLastAssigned, TimeToSend: user.TimeToSend, Studies: user.Studies}, nil
}


