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
	resp, err := client.CreateStudy(ctx, &pb.StudyMetaData{Name: "Flu Study", Description: "Flu study for winter season 2019"})
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
	fmt.Print(resp)
}

func TestCreateTrigger(t *testing.T)  {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudyClient(conn)
	conditions := []string{"age>25", "sick=yes", "weight>50"}
	action:= &pb.Action{Type: "time", Value: "3600000"}
	var actions []*pb.Action
	actions = append(actions, action)

	trigger := pb.Trigger{Condition: conditions, Action: actions, StudyID: "5d714e727ff561a1290d3df9"}
	resp, err := client.CreateTrigger(ctx, &pb.Trigger{Condition: trigger.Condition, StudyID: trigger.StudyID, Action: trigger.Action})
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
	fmt.Print(resp)
}

func TestStartStudy(t *testing.T)  {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudyClient(conn)

	resp, err := client.StartStudy(ctx, &pb.StudyMetaData{Id: "5d714e727ff561a1290d3df9"})
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
	fmt.Print(resp)
}

func TestUserSignUp(t *testing.T)  {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudyClient(conn)

	resp, err := client.UserSignUp(ctx, &pb.SignUpData{StudyID: "5d714e727ff561a1290d3df9"})
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
	fmt.Print(resp)
}

func TestCheckTriggers(t *testing.T){
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudyClient(conn)

	resp, err := client.CheckTriggers(ctx, &pb.Attributes{UserID: "5d7187a0737d658bb2e69cf5", Age: 26, Sick: "yes", Weight: 70})
	if err != nil {
		log.Fatalf("Error on Add: %v", err)
	}
	fmt.Println(resp)
}

func TestGetAssignedSurvey(t *testing.T)  {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudyClient(conn)

	resp, err := client.GetAssignedSurvey(ctx, &pb.AssignedSurvey{Id: "5d719dc601849cbd30a21405"})
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
	fmt.Print(resp)
}

func TestSubmitSurvey(t *testing.T)  {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudyClient(conn)

	qType := "matrix"

	textEn := "Tell us about your week"
	textDe := "Erzählen Sie uns von Ihrer Woche."
	text := make(map[string]string)
	text["english"] = textEn
	text["german"] = textDe

	matrixText1En := "How active have you been?"
	matrixText1De := "Wie aktiv waren Sie?"
	matrixText1 := make(map[string]string)
	matrixText1["english"] = matrixText1En
	matrixText1["german"] = matrixText1De

	options1En := pb.Question_Options{OptionText: []string{"very active", "active", "not active at all"}}
	options1De := pb.Question_Options{OptionText: []string{"sehr aktiv", "aktiv", "nicht aktiv"}}
	options1 := make(map[string] *pb.Question_Options)
	options1["english"] = &options1En
	options1["german"] = &options1De

	matrixOption1 := &pb.Question_MatrixOptions{Text: matrixText1, Options: options1, SelectedOption: "very active"}

	matrixText2En := "How are you feeling now?"
	matrixText2De := "Wie geht es Ihnen heute?"
	matrixText2 := make(map[string]string)
	matrixText2["english"] = matrixText2En
	matrixText2["german"] = matrixText2De

	options2En := pb.Question_Options{OptionText: []string{"good", "normal", "bad"}}
	options2De := pb.Question_Options{OptionText: []string{"gut", "normal", "schlecht"}}
	options2 := make(map[string] *pb.Question_Options)
	options2["english"] = &options2En
	options2["german"] = &options2De

	matrixOption2 := &pb.Question_MatrixOptions{Text: matrixText2, Options: options2, SelectedOption: "normal"}

	var matrixOptions []*pb.Question_MatrixOptions
	matrixOptions = append(matrixOptions, matrixOption1, matrixOption2)

	answerOptions := &pb.Question_AnswerOptions{MatrixOptions: matrixOptions}

	questionWithLanguage := &pb.Question_QuestionWithLanguage{Text: text, AnswerOptions : answerOptions}

	var question1 = &pb.Question{Id: "5d6928bee278a85b3f2c96d6",Type: qType, QuestionWithLanguage: questionWithLanguage}

	var questions []*pb.Question
	questions = append(questions, question1)

	var survey = &pb.SurveyData{Id: "5d692a61e278a85b3f2c96d7", Description: "This is a Survey", Type: "timely", Questions: questions}

	resp, err := client.SubmitSurvey(ctx, &pb.Answer{AssignmentID: "5d719dc601849cbd30a21405", Survey: survey})
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
	fmt.Print(resp)
}

func TestFinishStudy(t *testing.T){
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial localhost: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudyClient(conn)

	resp, err := client.FinishStudy(ctx, &pb.StudyMetaData{Id: "5d714e727ff561a1290d3df9"})
	if err != nil {
		log.Fatalf("Error on Add: %v", err)
	}
	fmt.Println(resp)
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