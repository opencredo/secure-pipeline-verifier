package cmd

import (
	"os"
	"github/secure-pipeline-verifier/app/config"
	"github/secure-pipeline-verifier/app/policies/common"
	"github/secure-pipeline-verifier/app/policies/github"
	"github/secure-pipeline-verifier/app/policies/gitlab"
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
