package gitlab

import (
	"fmt"
	"secure-pipeline-poc/app/clients/gitlab"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

type Controls struct {
	Api *gitlab.Api
}


func (c *Controls) SetClient(token string){
	c.Api = gitlab.NewApi(token)
}

func (c *Controls) ValidateC1(policyPath string, cfg *config.Config, sinceDate time.Time) {
	fmt.Println("------------------------------Control-1------------------------------")
	policy := common.UserAuthPolicy(policyPath)
	ciCommits, err := c.Api.Repo.GetChangesToCiCd(
		cfg.RepoInfoChecks.CiCdPath,
		cfg.Project.Owner + "/" + cfg.Project.Repo,
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

	signatureProtection := c.Api.GetProjectSignatureProtection(
		cfg.Project.Owner + "/" + cfg.Project.Repo ,
	)
	policy := common.SignatureProtectionPolicy(policyPath)
	policy.Process(cfg.Slack, common.GetObjectMap(signatureProtection))
}

func (c *Controls) ValidateC3(policyPath string, cfg *config.Config) {
	fmt.Println("------------------------------Control-3------------------------------")

	automationKeys, _ := c.Api.GetAutomationKeys(
		cfg.Project.Owner + "/" + cfg.Project.Repo ,
	)
	policy := common.KeyExpiryPolicy(policyPath)
	for _, item := range automationKeys {
		policy.Process(cfg.Slack, common.GetObjectMap(item))
	}
}

func (c *Controls) ValidateC4(policyPath string, cfg *config.Config) {
	fmt.Println("------------------------------Control-4------------------------------")
	automationKeys, _ := c.Api.GetAutomationKeys(
		cfg.Project.Owner + "/" + cfg.Project.Repo ,
	)

	policy := common.KeyReadOnlyPolicy(policyPath)
	for _, item := range automationKeys {
		policy.Process(cfg.Slack, common.GetObjectMap(item))
	}
}

