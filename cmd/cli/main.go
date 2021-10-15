package main

import (
	"fmt"
	"os"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"secure-pipeline-poc/app/policies/github"
	"secure-pipeline-poc/app/policies/gitlab"
	"time"
)

const (
	GitHubPlatform = "github"
	GitLabPlatform = "gitlab"
)


func PerformCheck (p *common.Platform){
	var envKey string
	switch p.Config.Project.Platform {
	case GitHubPlatform:
		envKey = config.GitHubToken
		p.Handler = &github.Handler{}
	case GitLabPlatform:
		envKey = config.GitLabToken
		p.Handler = &gitlab.Handler{}
	default:
		panic("Could not determine the platform!")
	}
	token := os.Getenv(envKey)
	p.Handler.SetClient(token)
	p.ValidatePolicies()
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "path/to/config/", "YYYY-MM-ddTHH:mm:ss.SSSZ")
		return
	}

	var cfg config.Config
	config.LoadConfig(os.Args[1], &cfg)
	config.LoadTrustedDataToJsonMap(os.Args[1], &cfg)

	sinceDate, err := time.Parse(time.RFC3339, os.Args[2])
	if err != nil {
		fmt.Println("Error " + err.Error() + " occurred while parsing date from " + os.Args[2])
		os.Exit(2)
	}

	platform := &common.Platform{
		Config: &cfg,
		SinceDate: sinceDate,
	}

	PerformCheck(platform)
}
