package main

import (

"context"
"encoding/json"
"fmt"
"github.com/aws/aws-lambda-go/events"
"github.com/aws/aws-lambda-go/lambda"
ac "github.com/aws/aws-sdk-go-v2/config"
lsvc "github.com/aws/aws-sdk-go-v2/service/lambda"
"github.com/aws/aws-sdk-go-v2/service/lambda/types"
"github.com/aws/aws-sdk-go/aws"
"net/url"
"os"
"strings"

)

const (
	targetLambda = "TARGET_LAMBDA"
)

type slackRequest struct {
	token       string
	teamID      string
	teamDomain  string
	channelID   string
	channelName string
	userID      string
	userName    string
	command     string
	text        string
	responseURL string
	triggerID   string
}

type PoliciesCheckEvent struct {
	Region   string `json:"region"`
	Bucket   string `json:"bucket"`
	RepoPath string `json:"configPath"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// valification token
	vals, _ := url.ParseQuery(request.Body)
	req := slackRequest{
		vals.Get("token"),
		vals.Get("team_id"),
		vals.Get("team_domain"),
		vals.Get("channel_id"),
		vals.Get("channel_name"),
		vals.Get("user_id"),
		vals.Get("user_name"),
		vals.Get("command"),
		vals.Get("text"),
		vals.Get("response_url"),
		vals.Get("trigger_id"),
	}

    // Slack arguments: '<param1> <param2> <param3>'
	args := strings.Fields(req.text)
 
	payload, err := json.Marshal(PoliciesCheckEvent{
        os.Getenv("AWS_REGION"),
		"secure-pipeline-bucket",
		args[0],
	})
	if err != nil {
		exitErrorf("Failed to load payload.", err.Error())
	}

	input := &lsvc.InvokeInput{
		FunctionName: aws.String(os.Getenv(targetLambda)),
		InvocationType: types.InvocationTypeEvent,
		Payload:      payload,
		LogType:      types.LogTypeTail,
	}
	
    cfg, err := ac.LoadDefaultConfig(ctx)
    if err != nil {
        exitErrorf("Failed to load default config.", err.Error())
    }
	client := lsvc.NewFromConfig(cfg)
	resp, err := client.Invoke(ctx, input)
	
	fmt.Printf("%v", payload)
	
	if err != nil {
		exitErrorf("Failed to invoke the Lambda function.", err.Error())
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%v", resp),
		StatusCode: 200,
	}, nil
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
