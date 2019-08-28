package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"labix.org/v2/mgo/bson"
	"log"
)

var studyCollection = client.Database("test").Collection("study")

type Study struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
	Name string
	Description string
	StartDate int64
	Status string
	Users []string
}

func createStudyDocument(doc Study) {
	insertResult, err := studyCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(insertResult)
}

func getAllStudies() []*Study{
	var studyResults []*Study
	if err != nil{
		log.Fatal(err)
	}

	findOptions := options.Find()
	cur, err := studyCollection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Study
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		studyResults = append(studyResults, &elem)
	}
	return studyResults
}

func getActiveStudies(status string) []*Study{
	var studyResults []*Study
	findOptions := options.Find()
	filter := bson.M{"status": status}

	cur, err := studyCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Study
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		studyResults = append(studyResults, &elem)
	}
	return studyResults
}


