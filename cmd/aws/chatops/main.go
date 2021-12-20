package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// valification token
	token := os.Getenv("VERIFICATION_TOKEN")
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

	if req.token != token {
		// invalid token
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Invalid token."),
			StatusCode: 401,
		}, nil
	}

	// ### Write your command logic ###

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%v", req),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}