package github

import (
	"fmt"
	x "github.com/google/go-github/v38/github"
	"secure-pipeline-poc/app/clients/github"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

type Controls struct {
	Client *x.Client
}


func (c *Controls) SetClient(token string){
	c.Client = github.NewClient(token)
}

func (c *Controls) ValidateC1(policyPath string, cfg *config.Config, sinceDate time.Time) {
	fmt.Println("------------------------------Control-1------------------------------")
	var policy = common.UserAuthPolicy(policyPath)

	ciCommits, err := github.GetChangesToCiCd(
		c.Client,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	if ciCommits != nil {
		for _, item := range ciCommits {
			policy.Process(cfg.Slack, common.GetObjectMap(item), cfg.RepoInfoChecks.TrustedData)
		}
		return
	}
	if err != nil {
		fmt.Printf("[Control 1: ERROR - performing control-1: %v]", err.Error())
		return
	}
	if ciCommits == nil {
		msg := fmt.Sprintf("{ \"control\": \"Control 1\", \"level\": \"WARNING\", \"msg\": \"No new commits since %v\"}", sinceDate)
		fmt.Println(msg)
		common.SendNotification(msg, cfg.Slack)
	}
}

func (c *Controls) ValidateC2(policyPath string, cfg *config.Config) {
	fmt.Println("------------------------------Control-2------------------------------")

	var policy = common.SignatureProtectionPolicy(policyPath)
	signatureProtection := github.GetBranchSignatureProtection(
		c.Client,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.ProtectedBranches,
	)
	for _, item := range signatureProtection {
		policy.Process(cfg.Slack, common.GetObjectMap(item))
	}}

func (c *Controls) ValidateC3(policyPath string, cfg *config.Config) {
	fmt.Println("------------------------------Control-3------------------------------")

	var policy = common.KeyExpiryPolicy(policyPath)
	automationKeysE, err := github.GetAutomationKeysExpiry(
		c.Client,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	if automationKeysE != nil {
		for _, item := range automationKeysE {
			policy.Process(cfg.Slack, common.GetObjectMap(item))
		}
	}
	if err != nil {
		fmt.Println("Error performing control-3: ", err)
	}
}

func (c *Controls) ValidateC4(policyPath string, cfg *config.Config) {
	fmt.Println("------------------------------Control-4------------------------------")

	var policy = common.KeyReadOnlyPolicy(policyPath)
	automationKeysRO, err := github.GetAutomationKeysPermissions(
		c.Client,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	if automationKeysRO != nil {
		for _, item := range automationKeysRO {
			policy.Process(cfg.Slack, common.GetObjectMap(item))
		}
	}
	if err != nil {
		fmt.Println("Error performing control-4: ", err)
	}
}
