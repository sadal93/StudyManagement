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
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"

	pb "Thesis-demo/api"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var clientOptions  = options.Client().ApplyURI("mongodb://localhost:27017")
var client, err = mongo.Connect(context.TODO(), clientOptions)
var collection = client.Database("test").Collection("additions")
var results []*Addition

type server struct{}

type Addition struct {
	Number1 int32
	Number2  int32
	SumResult int32
}

func getDocuments() []*Addition{

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
}

func (s *server) Add(ctx context.Context, in *pb.Input) (*pb.Output, error) {
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

}

func (s *server) MultipleSum(in *pb.Range, stream pb.Addition_MultipleSumServer) error {
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
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAdditionServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err )
	}
}
