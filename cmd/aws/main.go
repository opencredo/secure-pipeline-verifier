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
	"github.com/aws/aws-sdk-go/service/ssm"
	"io"
	"os"
	"path"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/github"
	"secure-pipeline-poc/app/policies/gitlab"
	"strings"
	"time"
)

const (
	GitHubPlatform = "github"
	GitLabPlatform = "gitlab"

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

func HandleRequest(ctx context.Context, event PoliciesCheckEvent) (string, error) {
	fmt.Printf("Running Policies Checks for Repo: %s \n", event.RepoPath)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(event.Region),
	})
	if err != nil {
		exitErrorf("Unable to create a new session %v", err)
	}

	var cfg config.Config
	loadConfig(event, sess, &cfg)
	lastRun := getLastRunParameterValue(sess)
	sinceDate, err := time.Parse(LastRunFormat, lastRun)

	if cfg.Project.Platform == GitHubPlatform {
		policiesObjList := collectPoliciesListFromS3(sess, event, GitHubPlatform)
		downloadPoliciesFromS3(sess, policiesObjList, event)
		var gitHubToken = os.Getenv(config.GitHubToken)
		github.ValidatePolicies(gitHubToken, &cfg, sinceDate)
	}
	if cfg.Project.Platform == GitLabPlatform {
		policiesObjList := collectPoliciesListFromS3(sess, event, GitLabPlatform)
		downloadPoliciesFromS3(sess, policiesObjList, event)
		var gitLabToken = os.Getenv(config.GitLabToken)
		gitlab.ValidatePolicies(gitLabToken, &cfg, sinceDate)
	}

	updateLastRunParameterValue(sess)
	return fmt.Sprintf("Check Complete for %s repo", event.RepoPath), nil
}

func loadConfig(event PoliciesCheckEvent, session *session.Session, cfg *config.Config) {
	svc := s3.New(session)
	configReadCloser := downloadConfigFromS3(svc, event.Bucket, event.RepoPath+"/"+config.ConfigsFileName)
	config.DecodeConfigToStruct(configReadCloser, cfg)

	trustedDataCloser := downloadConfigFromS3(svc, event.Bucket, event.RepoPath+"/"+config.TrustedDataFileName)
	config.DecodeTrustedDataToMap(trustedDataCloser, cfg)
}

func downloadConfigFromS3(svc *s3.S3, bucket string, item string) io.ReadCloser {
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

func collectPoliciesListFromS3(session *session.Session, event PoliciesCheckEvent, platform string) *s3.ListObjectsV2Output {
	svc := s3.New(session)

	policyObjects, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: aws.String(event.Bucket),
			Prefix: aws.String(event.RepoPath + PoliciesFolder + platform),
		},
	)
	if err != nil {
		exitErrorf("Unable to list items in bucket %q on folder %s, %v", event.Bucket, event.RepoPath+PoliciesFolder+platform, err)
	}

	return policyObjects
}

func downloadPoliciesFromS3(session *session.Session, policyObjects *s3.ListObjectsV2Output, event PoliciesCheckEvent) {
	downloader := s3manager.NewDownloader(session)

	fmt.Println("Downloading Policies found in S3:")
	for _, policyObject := range policyObjects.Contents {
		if strings.HasSuffix(*policyObject.Key, RegoExtension) {
			fmt.Println("Name:	", path.Base(*policyObject.Key))

			file, err := os.Create("/tmp/" + path.Base(*policyObject.Key))
			if err != nil {
				exitErrorf("Unable to open file %q, %v", path.Base(*policyObject.Key), err)
			}
			defer file.Close()

			numBytes, err := downloader.Download(file,
				&s3.GetObjectInput{
					Bucket: aws.String(event.Bucket),
					Key:    aws.String(*policyObject.Key),
				})
			if err != nil {
				exitErrorf("Unable to download item %q, %v", *policyObject.Key, err)
			}
			fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
		}
	}
}

func getLastRunParameterValue(session *session.Session) string {
	ssmsvc := ssm.New(session)
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(LastRunParameter),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		exitErrorf(err.Error())
	}

	value := *param.Parameter.Value
	fmt.Println("Retrieved value for the last performed checks: ", value)
	return value
}

func updateLastRunParameterValue(session *session.Session) {
	ssmsvc := ssm.New(session)
	_, err := ssmsvc.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(LastRunParameter),
		Value:     aws.String(time.Now().Format(LastRunFormat)),
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		exitErrorf(err.Error())
	}

	fmt.Println("Updated value for the last performed checks")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
