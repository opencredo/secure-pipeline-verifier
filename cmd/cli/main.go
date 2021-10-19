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


func PerformCheck (cfg *config.Config, sinceDate time.Time){
	var envKey string
	var Controls common.Controls
	switch cfg.Project.Platform {
	case GitHubPlatform:
		envKey = config.GitHubToken
		Controls = &github.Controls{}
	case GitLabPlatform:
		envKey = config.GitLabToken
		Controls = &gitlab.Controls{}
	default:
		panic("Could not determine the platform!")
	}
	input := &common.ValidateInput{
		Config:    cfg,
		Controls:  Controls,
		SinceDate: sinceDate,
		Token:     os.Getenv(envKey),
	}
	common.ValidatePolicies(input)
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

	PerformCheck(&cfg, sinceDate)
}
