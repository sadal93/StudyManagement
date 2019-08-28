package main

import (
	pb "StudyManagement/api"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestCreateStudy(t *testing.T)  {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudyClient(conn)
	resp, err := client.CreateStudy(ctx, &pb.StudyMetaData{Name: "Test Study", Description: "Testing study creation"})
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
	fmt.Print(resp)
}

func TestGetStudy(t *testing.T){
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudyClient(conn)

	resp, err := client.GetStudy(ctx, &pb.StudyMetaData{Id: "abcd"})
	if err != nil {
		log.Fatalf("Error on Add: %v", err)
	}
	fmt.Println(resp)
}