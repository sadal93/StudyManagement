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
	"log"
	"time"

	pb "Thesis-demo/api"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
)

func createUser( c pb.UserClient, user pb.UserMetaData){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateUser(ctx, &pb.UserMetaData{UserID: user.UserID, TimeLastAssigned: user.TimeLastAssigned, TimeToSend: user.TimeToSend, Role: user.Role})
	if err != nil {
		log.Fatalf("could not insert: %v", err)
	}

	log.Print("User = ", r.UserID)

}

func assignUserToStudy(c pb.StudyClient, studyUser pb.UserAssignment)  {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AssignUserToStudy(ctx, &pb.UserAssignment{StudyID: studyUser.StudyID, UserID: studyUser.UserID })
	if err != nil {
		log.Fatalf("could not assign: %v", err)
	}

	log.Print("User = ", r.Users)
	}

func createStudy( c pb.StudyClient, study pb.StudyMetaData){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateStudy(ctx, &pb.StudyMetaData{StudyID: study.StudyID, Name: study.Name, Description: study.Description, StartDate: study.StartDate, Status: study.Status, Users: study.Users})
	if err != nil {
		log.Fatalf("could not insert: %v", err)
	}

	log.Print("Study = ", r.StudyID)

}

func triggerCheck( c pb.CheckClient, user pb.UserMetaData){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//fmt.Print("Hello")
	r, err := c.CheckTrigger(ctx, &pb.UserMetaData{UserID: user.UserID, TimeLastAssigned: user.TimeLastAssigned, TimeToSend: user.TimeToSend, Role: user.Role})
	//fmt.Print("World")
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

/*func getSum( c pb.AdditionClient, first int32, second int32){

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
}*/

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//c := pb.NewAdditionClient(conn)
	//d := pb.NewCheckClient(conn)
	//e := pb.NewUserClient(conn)
	f := pb.NewStudyClient(conn)
	//rules := []string{"age=20", "sick=false", "weight=60"}
	//action := "survey1"

	//rule := pb.Rule{Condition: rules, Action: action}

	/*userID := "1"
	var timeLastAssigned int64 = 26
	var timeToSend int64 = 7000
	role := "participant"

	user := pb.UserMetaData{UserID: userID, TimeLastAssigned: timeLastAssigned, TimeToSend:timeToSend, Role: role}*/

	/*studyID := "1"
	name:= "Flu Study"
	description := "Study about flu"
	var startDate int64 = time.Now().UnixNano() / 1000000
	status := "Active"
	users := []string{""}

	study := pb.StudyMetaData{StudyID: studyID, Name: name, Description: description, StartDate: startDate, Status: status, Users: users}
	createStudy(f, study)*/

	userID := "1"
	studyID := "1"
	userAssignment := pb.UserAssignment{StudyID: studyID, UserID: userID}

	assignUserToStudy(f, userAssignment)

	//createUser(e, user)
	//for {
		//triggerCheck(d, user)
	//	time.Sleep(time.Second * 10)
	//}
	//createUser(e, user)
	//createRule(d, rule)
	//getSum(c, 2, 4)
	//getStreamSum(c, 1, 7)



}
