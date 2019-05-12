/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

// Package main implements a server for Greeter service.
package main

import (
	"fmt"
	"reflect"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"strconv"
	"strings"

	pb "Thesis-demo/api"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var clientOptions  = options.Client().ApplyURI("mongodb://localhost:27017")
var client, err = mongo.Connect(context.TODO(), clientOptions)
//var collection = client.Database("test").Collection("additions")
var triggerCollection = client.Database("test").Collection("triggers")
var userCollection = client.Database("test").Collection("users")
//var results []*Addition
var results []*Check

type server struct{}

/*type Addition struct {
	Number1 int32
	Number2  int32
	SumResult int32
}*/

type Check struct {
	Condition[] string
	Action string
}

type User struct {
	name string
	age int32
	sick string
	weight int32
}

func contains(condtions []bool, search bool) bool {
	for _, value := range condtions {
		if value == search {
			return true
		}
	}
	return false
}

/*func getDocuments() []*Addition{

	if err != nil{
		log.Fatal(err)
	}

	findOptions := options.Find()
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Addition
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}
	return results
}*/

func createTriggerDocument(doc Check) {
	insertResult, err := triggerCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(insertResult)
}

func createUserDocument(doc User) {
	insertResult, err := userCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(insertResult)
}

func getAllTriggers() []*Check{
	if err != nil{
		log.Fatal(err)
	}

	findOptions := options.Find()
	cur, err := triggerCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Check
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}
	return results
}

func (s *server)  User(ctx context.Context, user *pb.Attributes) (*pb.Attributes, error) {

	userDoc:= User{user.Name, user.Age, user.Sick, user.Weight}
	createUserDocument(userDoc)

	log.Printf("User Created: %v", user.Name)
	return &pb.Attributes{Name: user.Name, Age: user.Age, Sick: user.Sick, Weight: user.Weight}, nil
}

func (s *server)  Trigger(ctx context.Context, rule *pb.Rule) (*pb.Rule, error) {

	trigger:= Check{rule.Condition, rule.Action}
	createTriggerDocument(trigger)

	log.Printf("Trigger Created: %v", trigger.Condition)
	return &pb.Rule{Condition: trigger.Condition, Action: trigger.Action}, nil
}

func (s *server)  CheckTrigger(ctx context.Context, user *pb.Attributes) (*pb.Rule, error) {
	documents := getAllTriggers()
	fmt.Println("USER: " , user)
	var action string

	userObj := User{user.Name, user.Age, user.Sick, user.Weight}

	for _,document := range documents {
		checks := document.Condition
		fmt.Println("Document: ", document)
		var conditions []bool

		for  _,check := range checks{
			condition := false
			var operator string
			var attribute string
			var value int64
			var valueString string
			 if strings.Contains(check, "<"){
			 	condition := strings.Split(check, "<")
			 	attribute = condition[0]
				 if i, err := strconv.Atoi(condition[1]); err == nil {
					 fmt.Println(" i: ", i)
					 value = int64(i)
					 fmt.Println(" value: ", value)
				 } else {
					 valueString = condition[1]
					 fmt.Println(" value: ", valueString)
				 }
				 fmt.Println(" attribute: ", attribute)

				 operator = "<"
				 fmt.Println(" operator: ", operator)

			 } else if strings.Contains(check, ">"){
				condition := strings.Split(check, ">")
				attribute = condition[0]
				 if i, err := strconv.Atoi(condition[1]); err == nil {
					 fmt.Println(" i = ", i)
					 value = int64(i)
					 fmt.Println(" value: ", value)
				 } else {
					 valueString = condition[1]
				 	fmt.Println(" value: ", valueString)
				 }

				 fmt.Println(" attribute: ", attribute)

				operator = ">"
				fmt.Println(" operator: ", operator)

			} else {
				 condition := strings.Split(check, "=")
				 attribute = condition[0]
				 if i, err := strconv.Atoi(condition[1]); err == nil {
					 fmt.Println(" i: ", i)
					 value = int64(i)
					 fmt.Println(" value: ", value)
				 } else {
					 valueString = condition[1]
					 fmt.Println(" value: ", valueString)
				 }

				 fmt.Println(" attribute: ", attribute)

				 operator = "="
				 fmt.Println(" operator: ", operator)
			 }

			switch operator {
			case ">":
				rv := reflect.ValueOf(userObj)
				//rv = rv.Elem()
				fmt.Println("  rv: ", rv)
				fmt.Println("  field by name: ", rv.FieldByName(attribute))
				if rv.FieldByName(attribute).Int() > value{
					condition = true
				}
				fmt.Println(" CONDITION ", condition)
				conditions = append(conditions, condition)

			case "<":
				rv := reflect.ValueOf(userObj)
				//rv = rv.Elem()
				fmt.Println("  rv: ", rv)
				fmt.Println("  field by name: ", rv.FieldByName(attribute))
				if rv.FieldByName(attribute).Int() < value{
					condition = true
				}
				fmt.Println(" CONDITION ", condition)
				conditions = append(conditions, condition)

			case "=":
				rv := reflect.ValueOf(userObj)
				//rv = rv.Elem()
				fmt.Println("  rv: ", rv)
				fmt.Println("  field by name: ", rv.FieldByName(attribute))
				if  rv.FieldByName(attribute).String() == valueString{
					condition = true
				}
				fmt.Println(" CONDITION ", condition)
				conditions = append(conditions, condition)
			}

		}
		fmt.Println("CONDITIONS", conditions)
		fmt.Println("")

		if !contains(conditions, false){
			action = document.Action
		}
	}


	return &pb.Rule{ Action: action}, nil
}

/*func (s *server) Add(ctx context.Context, in *pb.Input) (*pb.Output, error) {
	log.Printf("Received First: %v", in.First)
	log.Printf("Received Second: %v", in.Second)
	var a = in.First + in.Second
	documents := getDocuments()
	var sum int32 = 0
	for _,document := range documents {
		if document.Number1 == in.First && document.Number2 == in.Second {
			sum = document.SumResult
		}
	}
	if sum == 0 {
		log.Printf("Sum from Server: %v", a)
		return &pb.Output{Result: a}, nil
	} else {
		log.Printf("Sum from DB: %v", sum)
		return &pb.Output{Result: sum}, nil
	}

}*/

/*func (s *server) MultipleSum(in *pb.Range, stream pb.Addition_MultipleSumServer) error {
	log.Printf("Start Range: %v", in.Begin)
	log.Printf("End Range: %v", in.End)
	var sum int32 = 0
	for i := in.Begin; i <= in.End; i++ {
		sum += i
		err := stream.Send(&pb.Output{Result:sum})
		if err != nil {
			return err
		}

	}
	return nil
}*/

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	//pb.RegisterAdditionServer(s, &server{})
	pb.RegisterCheckServer(s, &server{})
	pb.RegisterGetUserServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err )
	}
}
