package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	lsvc "github.com/aws/aws-sdk-go/service/lambda"
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

// Commands coming from Slack comes in one string
type slackParameters struct {
	Repo string
	// TODO: We should be able to pass LastRun via slack
	LastRun string
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

	args := strings.Fields(req.command)
	params := slackParameters{
		Repo: args[0],
	}

	newSession := session.Must(session.NewSession())
	svc := lsvc.New(newSession)

	body, err := json.Marshal(map[string]interface{}{
		"region":     "eu-west-2",
		"bucket":     "secure-pipeline-bucket",
		"configPath": params.Repo,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	input := &lsvc.InvokeInput{
		FunctionName: aws.String("secure-pipeline"),
		Payload:      body,
	}
	invoke, err := svc.Invoke(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%v", invoke.Payload),
		StatusCode: 200,
	}, nil
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
