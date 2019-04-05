package tests

import (
	pb "Thesis-demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"testing"
)

func TestAddition(t *testing.T)  {

	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewAdditionClient(conn)
	resp, err := client.Add(ctx, &pb.Input{First:2, Second: 4})
	if err != nil {
		t.Fatalf("Sum failed: %v", err)
	}
	if resp.Result != 6{
		t.Errorf("Sum was incorrect, got: %d, want: %d.", resp.Result, 6)
	}
}

func TestMultipleAddition(t *testing.T){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewAdditionClient(conn)

	stream, err := client.MultipleSum(context.Background(), &pb.Range{Begin: 1, End: 7})
	if err != nil {
		log.Fatalf("Error on Add: %v", err)
	}
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Add(_) = _, %v", client, err)
		}
	}
}