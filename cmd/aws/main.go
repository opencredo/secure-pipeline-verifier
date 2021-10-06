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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

const (
	ConfigFileName = "config.yaml"
	TrustedDataFileName = "trusted-data.json"
)

type PoliciesCheckEvent struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
	Org string `json:"org"`
	Repo string `json:"repo"`
}

func HandleRequest(ctx context.Context, policiesCheckEvent PoliciesCheckEvent) (string, error) {
	repoPath := policiesCheckEvent.Org +"/"+ policiesCheckEvent.Repo
	fmt.Printf("Running Policies Checks for Repo: %s \n", repoPath)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(policiesCheckEvent.Region),
	})
	if err != nil {
		exitErrorf("Unable to create a new session %v", err)
	}

	downloader := s3manager.NewDownloader(sess)

	configFile, err := os.Create("/tmp/"+ConfigFileName)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", ConfigFileName, err)
	}
	downloadFileFromS3(downloader, configFile, policiesCheckEvent.Bucket, repoPath+"/"+ConfigFileName)

	trustedDataFile, err := os.Create("/tmp/"+TrustedDataFileName)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", TrustedDataFileName, err)
	}
	downloadFileFromS3(downloader, trustedDataFile, policiesCheckEvent.Bucket, repoPath+"/"+TrustedDataFileName)

	return fmt.Sprintf("Check Complete for %s repo", repoPath), nil
}

func main() {
	lambda.Start(HandleRequest)
}

func downloadFileFromS3(downloader *s3manager.Downloader, file *os.File, bucket string, item string) {
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}



