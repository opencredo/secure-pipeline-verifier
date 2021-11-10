package main

// Build Command:
// $ GOOS=linux go build -o main main.go

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	ac "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"io"
	"os"
	"path"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/cmd"
	"strings"
	"time"
)

const (
	S3PoliciesFolder = "/policies/"
	LambdaPoliciesFolder = "/tmp/"
	RegoExtension    = ".rego"

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
	fmt.Printf("Running Policies Checks for Controls: %s \n", event.RepoPath)

	awsCfg, err := ac.LoadDefaultConfig(ctx, ac.WithRegion(event.Region))

	if err != nil {
		exitErrorf("Failed loading AWS config, %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	var cfg config.Config
	loadConfig(ctx, event, s3Client, &cfg)

	ssmClient := ssm.NewFromConfig(awsCfg)
	lastRun := getLastRunParameterValue(ctx, ssmClient)
	sinceDate, err := time.Parse(LastRunFormat, lastRun)
	if err != nil {
		exitErrorf("Unable to read the date-time of last run %v", err)
	}

	policiesObjList := collectPoliciesListFromS3(ctx, s3Client, event)
	downloadPoliciesFromS3(ctx, s3Client, policiesObjList)
	cmd.PerformCheck(&cfg, sinceDate)

	updateLastRunParameterValue(ctx, ssmClient)
	return fmt.Sprintf("Check Complete for %s repo", event.RepoPath), nil
}

func loadConfig(ctx context.Context, event PoliciesCheckEvent, client *s3.Client, cfg *config.Config) {
	configReadCloser := downloadConfigFromS3(ctx, client, event.Bucket, event.RepoPath+"/"+config.ConfigsFileName)
	config.DecodeConfigToStruct(configReadCloser, cfg)
	updatePoliciesPath(cfg.RepoInfoChecks.Policies)
	trustedDataCloser := downloadConfigFromS3(ctx, client, event.Bucket, event.RepoPath+"/"+config.TrustedDataFileName)
	config.DecodeTrustedDataToMap(trustedDataCloser, cfg)
}

func downloadConfigFromS3(ctx context.Context, client *s3.Client, bucket string, item string) io.ReadCloser {
	resultInput := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &item,
	}

	result, err := client.GetObject(ctx, resultInput)
	if err != nil {
		exitErrorf("Unable to retrieve an S3 object: %v/%v", bucket, item, err.Error())
	}

	return result.Body
}

func collectPoliciesListFromS3(ctx context.Context, client *s3.Client, event PoliciesCheckEvent) *s3.ListObjectsV2Output {
	prefix := event.RepoPath + S3PoliciesFolder
	policyObjects, err := client.ListObjectsV2(ctx,
		&s3.ListObjectsV2Input{
			Bucket: &event.Bucket,
			Prefix: &prefix,
		},
	)
	if err != nil {
		exitErrorf("Unable to list items in bucket %q on folder %s, %v", event.Bucket, event.RepoPath+S3PoliciesFolder, err)
	}

	return policyObjects
}

func downloadPoliciesFromS3(ctx context.Context, client *s3.Client, policyObjects *s3.ListObjectsV2Output) {
	fmt.Println("Downloading Policies found in S3:")
	for _, policyObject := range policyObjects.Contents {
		if strings.HasSuffix(*policyObject.Key, RegoExtension) {
			fmt.Println("Name:	", path.Base(*policyObject.Key))

			file, err := os.Create(LambdaPoliciesFolder + path.Base(*policyObject.Key))
			if err != nil {
				exitErrorf("Unable to open file %q, %v", path.Base(*policyObject.Key), err)
			}
			defer file.Close()

			objectOut, err := client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: policyObjects.Name,
				Key:    policyObject.Key,
			})
			if err != nil {
				exitErrorf("Unable to download item %q, %v", *policyObject.Key, err)
			}

			numBytes, err := io.Copy(file, objectOut.Body)
			if err != nil {
				exitErrorf("Save item %q in the storage, %v", *policyObject.Key, err)
			}

			fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
		}
	}
}

func updatePoliciesPath(policies []config.Policies) {
	for i := 0; i < len(policies); i++ {
		policies[i].Path = LambdaPoliciesFolder + policies[i].Path
	}
}

func getLastRunParameterValue(ctx context.Context, client *ssm.Client) string {
	var name = LastRunParameter
	param, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: false,
	})
	if err != nil {
		exitErrorf(err.Error())
	}

	value := *param.Parameter.Value
	fmt.Println("Retrieved value for the last performed checks: ", value)
	return value
}

func updateLastRunParameterValue(ctx context.Context, client *ssm.Client) {
	var name = LastRunParameter
	var value = time.Now().Format(LastRunFormat)
	param, err := client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      &name,
		Value:     &value,
		Overwrite: true,
	})
	if err != nil {
		exitErrorf(err.Error())
	}

	fmt.Println("Updated value for the last performed checks ", param)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
