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
	GitHubPlatform = "github"
	GitLabPlatform = "gitlab"
)

func main(){

	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "path-to-config.yaml", "YYYY-MM-ddTHH:mm:ss.SSSZ")
		return
	}

	var cfg config.Config
	config.LoadConfig(os.Args[1], &cfg)

	sinceDate, err := time.Parse(time.RFC3339, os.Args[2])
	if err != nil {
		fmt.Println("Error " + err.Error() + " occurred while parsing date from " + os.Args[2])
		os.Exit(2)
	}

	if cfg.Project.Platform == GitHubPlatform {
		var gitHubToken = os.Getenv(config.GitHubToken)
		github.ValidatePolicies(gitHubToken, &cfg, sinceDate)
	}
	if cfg.Project.Platform == GitLabPlatform {
		var gitLabToken = os.Getenv(config.GitLabToken)
		gitlab.ValidatePolicies(gitLabToken, &cfg, sinceDate)
	}
}
