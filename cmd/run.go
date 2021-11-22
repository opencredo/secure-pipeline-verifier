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
	var controls common.Controls
	switch cfg.Project.Platform {
	case config.GitHubPlatform:
		controls = &github.Controls{}
	case config.GitLabPlatform:
		controls = &gitlab.Controls{}
	default:
		panic("Could not determine the platform!")
	}
	input := &common.ValidateInput{
		Config:    cfg,
		Controls:  controls,
		SinceDate: sinceDate,
		Token:     os.Getenv(config.RepoToken),
	}
	common.ValidatePolicies(input)
}
