package cmd

import (
	"os"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"secure-pipeline-poc/app/policies/github"
	"secure-pipeline-poc/app/policies/gitlab"
	"time"
)

func PerformCheck(cfg *config.Config, sinceDate time.Time) {
	var envKey string
	var controls common.Controls
	switch cfg.Project.Platform {
	case config.GitHubPlatform:
		envKey = config.GitHubToken
		controls = &github.Controls{}
	case config.GitLabPlatform:
		envKey = config.GitLabToken
		controls = &gitlab.Controls{}
	default:
		panic("Could not determine the platform!")
	}
	input := &common.ValidateInput{
		Config:    cfg,
		Controls:  controls,
		SinceDate: sinceDate,
		Token:     os.Getenv(envKey),
	}
	common.ValidatePolicies(input)
}
