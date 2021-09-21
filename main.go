package main

import (
	"fmt"
	"os"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/github"
	"secure-pipeline-poc/app/policies/gitlab"
	"time"
)

const (
	GitHubToken = "GITHUB_TOKEN"
	GitLabToken = "GITLAB_TOKEN"
)


var token = os.Getenv(GitHubToken)

func main()  {

	if len(os.Args) != 4 {
		fmt.Println("Usage:", os.Args[0], "path-to-config.yaml", "YYYY-MM-ddTHH:mm:ss.SSSZ" , "github/gitlab")
		return
	}

	var cfg config.Config
	config.LoadConfig(os.Args[1], &cfg)

	sinceDate, err := time.Parse(time.RFC3339, os.Args[2])
	if err != nil {
		fmt.Println("Error "+ err.Error() +" occurred while parsing date from "+ os.Args[2])
		os.Exit(2)
	}

	if os.Args[3] == "github" {
		var token = os.Getenv(GitHubToken)
		github.ValidatePolicies(token, cfg, sinceDate)
	}

	if os.Args[3] == "gitlab" {
		var token = os.Getenv(GitLabToken)
		gitlab.ValidatePolicies(token)
	}

}

