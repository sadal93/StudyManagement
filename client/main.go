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
	r, err := c.CreateStudy(ctx, &pb.StudyMetaData{Name: study.Name, Description: study.Description, StartDate: study.StartDate, Status: study.Status, Users: study.Users})
	if err != nil {
		log.Fatalf("could not insert: %v", err)
	}

	log.Print("Study = ", r.Name)

}

func getStudies( c pb.StudyClient){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetAll(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not insert: %v", err)
	}
	log.Print("Studies = ", r)
}

func updateStudy(c pb.StudyClient, study pb.StudyMetaData){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UpdateStudy(ctx, &pb.StudyMetaData{Id: study.Id, Name: study.Name, Description: study.Description, Status: study.Status})
	if err != nil {
		log.Fatalf("could not update: %v", err)
	}

	log.Print("Study = ", r.Name)
}

func deleteStudy(c pb.StudyClient, study pb.StudyID){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.DeleteStudy(ctx, &pb.StudyID{StudyID: study.StudyID})
	if err != nil {
		log.Fatalf("could not delete: %v", err)
	}

	log.Print("Study = ", r)
}

func getUsers(c pb.StudyClient, studyID pb.StudyID){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetUsers(ctx, &pb.StudyID{StudyID: studyID.StudyID})
	if err != nil {
		log.Fatalf("could not retrieve: %v", err)
	}

	log.Print("Users = ", r.Users)

}

func getStudy(c pb.StudyClient, studyID pb.StudyID){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetStudy(ctx, &pb.StudyID{StudyID: studyID.StudyID})
	if err != nil {
		log.Fatalf("could not retrieve: %v", err)
	}

	log.Print("Study = ", r)

}

func checkTrigger( c pb.TriggerClient, attributes pb.Attributes){

	stream, err := c.CheckTrigger(context.Background(), &pb.Attributes{Age: attributes.Age, Sick: attributes.Sick, Weight: attributes.Weight})

	if err != nil {
		log.Fatalf("Error Occured: %v", err)
	}

	for {
		strm, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Add(_) = _, %v", c, err)
		}
		log.Print("Action: ", strm)
	}


}

func createTrigger( c pb.TriggerClient, trigger pb.CreatedTrigger){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateTrigger(ctx, &pb.CreatedTrigger{Condition: trigger.Condition, StudyID: trigger.StudyID, Action: trigger.Action})
	if err != nil {
		log.Fatalf("could not insert: %v", err)
	}

	log.Print("Trigger = ", r.Condition)

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
	//d := pb.NewTriggerClient(conn)
	//e := pb.NewUserClient(conn)
	f := pb.NewStudyClient(conn)
	/*conditions := []string{"age>25", "sick=yes", "weight>50"}
	action1 := &pb.Action{Type: "survey", Value: "1"}
	action2:= &pb.Action{Type: "time", Value: "3600000"}
	var actions []*pb.Action
	actions = append(actions, action1, action2)

	createdTrigger := pb.CreatedTrigger{Condition: conditions, Action: actions, StudyID: "5d0b6b28678629f9b50baa02"}
	createTrigger(d, createdTrigger)*/

	/*attributes := pb.Attributes{Age: 26, Sick: "yes", Weight: 70}
	checkTrigger(d, attributes)*/
	//action := "survey1"

	//rule := pb.Rule{Condition: rules, Action: action}

	/*userID := "1"
	var timeLastAssigned int64 = 26
	var timeToSend int64 = 7000
	role := "participant"

	user := pb.UserMetaData{UserID: userID, TimeLastAssigned: timeLastAssigned, TimeToSend:timeToSend, Role: role}*/

	/*id := "5d0c381c58103882ca0fdda4"
	studyToDelete := pb.StudyID{StudyID: id}
	deleteStudy(f, studyToDelete)*/
	//id := "5d0b6b28678629f9b50baa02"
	/*name:= "Headache Study"
	description := "Study about Headache"
	//var startDate int64 = time.Now().UnixNano() / 1000000
	status := "Active"
	//users := []string{""}

	study1 := pb.StudyMetaData{Name: name, Description: description, Status: status}
	createStudy(f, study1)*/
	//updateStudy(f, study1)


	/*userID := "1"
	studyID := "5d0b6b28678629f9b50baa02"
	userAssignment := pb.UserAssignment{StudyID: studyID, UserID: userID}

	assignUserToStudy(f, userAssignment)*/

	studyID := "5d0b6b28678629f9b50baa02"
	study2 := pb.StudyID{StudyID: studyID}

	getStudy(f, study2)

	getUsers(f, study2)

	//getStudies(f)

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
