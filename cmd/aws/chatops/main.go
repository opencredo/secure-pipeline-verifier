package main

// Build Command:
// $ GOOS=linux go build -o main main.go

import (
    "context"
    "fmt"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "os"
    "time"
)

const (
    PoliciesFolder = "/policies/"
    RegoExtension  = ".rego"

    LastRunParameter = "/Lambda/SecurePipelines/last_run"
    LastRunFormat    = time.RFC3339
)

type PoliciesCheckEvent struct {
    Region   string `json:"region"`
    Bucket   string `json:"bucket"`
    RepoPath string `json:"configPath"`
}

func main() {
    lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (string, error) {
    body := request.Body
    switch request.HTTPMethod {
    case "POST":

    }
}

func post() {

}

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}
