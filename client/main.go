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
	"io"
	"log"
	"time"

	pb "Thesis-demo/api"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
)

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

	c := pb.NewAdditionClient(conn)

	getSum(c, 2, 4)
	getStreamSum(c, 1, 7)



}
