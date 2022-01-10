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
	Branch   string `json:"branch,omitempty"`
}

type SlackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*SlackResponse, error) {

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

	if len(args) < 1 {
		exitErrorf("At least 1 argument required from Slack.")
	}

	data := &PoliciesCheckEvent{
		os.Getenv("AWS_REGION"),
		"secure-pipeline-bucket",
		args[0],
		"",
	}

	if len(args) > 1 {
		data.Branch = args[1]
	}

	payload, err := json.Marshal(data)
	if err != nil {
		exitErrorf("Failed to load payload.", err.Error())
	}

	input := &lsvc.InvokeInput{
		FunctionName:   aws.String(os.Getenv(targetLambda)),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
		LogType:        types.LogTypeTail,
	}

	// Call another Lambda
	cfg, err := ac.LoadDefaultConfig(ctx)
	if err != nil {
		exitErrorf("Failed to load default config.", err.Error())
	}
	client := lsvc.NewFromConfig(cfg)
	_, err = client.Invoke(ctx, input)
	if err != nil {
		exitErrorf(fmt.Sprintf("Failed to invoke the Lambda function: %v.", input.FunctionName), err.Error())
	}

	return &SlackResponse{
		ResponseType: "ephemeral",
		Text: fmt.Sprintf("Triggered a Secure Pipeline check for: %v, branch: %v",
			data.RepoPath, data.Branch),
	}, nil
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
