package main

// Build Command:
// $ GOOS=linux go build -o main main.go

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/github"
	"secure-pipeline-poc/app/policies/gitlab"
	"time"
)

const (
	GitHubPlatform = "github"
	GitLabPlatform = "gitlab"
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

	awsCfg := aws.Config{
		Region: policiesCheckEvent.Region,
	}
	sess, err := session.NewSession(&awsCfg)
	if err != nil {
		exitErrorf("Unable to create a new session %v", err)
	}

	var cfg config.Config
	loadConfig(policiesCheckEvent, sess, &cfg)

	// TODO fix this - hard-coding date for now
	sinceDate, err := time.Parse(time.RFC3339, "2020-01-01T09:00:00.000Z")
	if err != nil {
		fmt.Println("Error " + err.Error() + " occurred while parsing date from " + "2020-01-01T09:00:00.000Z")
		exitErrorf(err.Error())
	}

	if cfg.Project.Platform == GitHubPlatform {
		var gitHubToken = os.Getenv(config.GitHubToken)
		github.ValidatePolicies(gitHubToken, &cfg, sinceDate)
	}
	if cfg.Project.Platform == GitLabPlatform {
		var gitLabToken = os.Getenv(config.GitLabToken)
		gitlab.ValidatePolicies(gitLabToken, &cfg, sinceDate)
	}
	ProcessLastRun(ctx, awsCfg)
	return fmt.Sprintf("Check Complete for %s repo", repoPath), nil
}

func loadConfig(event PoliciesCheckEvent, session *session.Session, cfg *config.Config) {
	svc := s3.New(session)
	repoPath := event.Org + "/" + event.Repo
	configReadCloser := downloadFileFromS3(svc, event.Bucket, repoPath+"/"+config.ConfigsFileName)
	config.DecodeConfigToStruct(configReadCloser, cfg)

	trustedDataCloser := downloadFileFromS3(svc, event.Bucket, repoPath+"/"+config.TrustedDataFileName)
	config.DecodeTrustedDataToMap(trustedDataCloser, cfg)
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
