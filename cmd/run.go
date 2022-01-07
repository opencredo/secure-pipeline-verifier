package cmd

import (
	"os"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"secure-pipeline-poc/app/policies/github"
	"secure-pipeline-poc/app/policies/gitlab"
	"time"
)

func PerformCheck(cfg *config.Config, sinceDate time.Time, branch string) {
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
		Branch:    branch,
		Config:    cfg,
		Controls:  controls,
		SinceDate: sinceDate,
		Token:     os.Getenv(config.RepoToken),
	}
	common.ValidatePolicies(input)
}
