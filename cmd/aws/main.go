package main

// Build Command:
// $ GOOS=linux go build -o main main.go

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"os"
	"secure-pipeline-poc/app/config"
)

const (
	ConfigFileName      = "config.yaml"
	TrustedDataFileName = "trusted-data.json"
)

type PoliciesCheckEvent struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
	Org    string `json:"org"`
	Repo   string `json:"repo"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, policiesCheckEvent PoliciesCheckEvent) (string, error) {
	repoPath := policiesCheckEvent.Org + "/" + policiesCheckEvent.Repo
	fmt.Printf("Running Policies Checks for Repo: %s \n", repoPath)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(policiesCheckEvent.Region),
	})
	if err != nil {
		exitErrorf("Unable to create a new session %v", err)
	}

	svc := s3.New(sess)
	configReadCloser := downloadFileFromS3(svc, policiesCheckEvent.Bucket, repoPath+"/"+ConfigFileName)
	var cfg config.Config
	config.DecodeConfig(configReadCloser, &cfg)
	fmt.Println("Controls To Run: ", cfg.RepoInfoChecks.ControlsToRun)

	downloadFileFromS3(svc, policiesCheckEvent.Bucket, repoPath+"/"+TrustedDataFileName)

	return fmt.Sprintf("Check Complete for %s repo", repoPath), nil
}

func downloadFileFromS3(svc *s3.S3, bucket string, item string) io.ReadCloser {
	resultInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	}

	result, err := svc.GetObject(resultInput)
	if err != nil {
		exitErrorf(err.Error())
	}

	return result.Body
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
