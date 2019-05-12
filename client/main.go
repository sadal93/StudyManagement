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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "Thesis-demo/api"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
)

func createUser( c pb.GetUserClient, user pb.Attributes){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.User(ctx, &pb.Attributes{Name: user.Name, Age: user.Age, Sick: user.Sick, Weight: user.Weight})
	if err != nil {
		log.Fatalf("could not insert: %v", err)
	}

	log.Print("User = ", r.Name)

}

func triggerCheck( c pb.CheckClient, user pb.Attributes){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Print("Hello")
	r, err := c.CheckTrigger(ctx, &pb.Attributes{Name: user.Name, Age: user.Age, Sick: user.Sick, Weight: user.Weight})
	fmt.Print("World")
	if err != nil {
		log.Fatalf("Error Occured: %v", err)
	}

	log.Print("Action = ", r.Action)

}

func createRule( c pb.CheckClient, rule pb.Rule){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Trigger(ctx, &pb.Rule{Condition: rule.Condition, Action: rule.Action})
	if err != nil {
		log.Fatalf("could not insert: %v", err)
	}

	log.Print("Rule = ", r.Condition)

}

func getSum( c pb.AdditionClient, first int32, second int32){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Add(ctx, &pb.Input{First: first, Second: second})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}

	log.Print("Sum = ", r.Result)

}

func getStreamSum(c pb.AdditionClient,  begin int32, end int32)  {

	stream, err := c.MultipleSum(context.Background(), &pb.Range{Begin: begin, End: end})
	if err != nil {
		log.Fatalf("Error on Add: %v", err)
	}
	for {
		sum, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Add(_) = _, %v", c, err)
		}
		log.Print("Sum: ", sum)
	}
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//c := pb.NewAdditionClient(conn)
	d := pb.NewCheckClient(conn)
	//e := pb.NewGetUserClient(conn)

	//rules := []string{"age=20", "sick=false", "weight=60"}
	//action := "survey1"

	//rule := pb.Rule{Condition: rules, Action: action}

	name := "John"
	var age int32 = 26
	sick := "yes"
	var weight int32 = 77

	user := pb.Attributes{Name: name, Age: age, Sick: sick, Weight: weight}

	//createUser(e, user)

	triggerCheck(d, user)

	//createRule(d, rule)
	//getSum(c, 2, 4)
	//getStreamSum(c, 1, 7)



}
