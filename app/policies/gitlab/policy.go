package gitlab

import (
	"fmt"
	"secure-pipeline-poc/app/clients/gitlab"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

func ValidatePolicies(token string, cfg *config.Config, sinceDate time.Time) {
	api := gitlab.NewApi(token, cfg)

	for _, policy := range cfg.RepoInfoChecks.Policies {
		switch policy.Control {
		case config.Control1:
			if policy.Enabled {
				ValidateC1(api, cfg, policy.Path, sinceDate)
			}
		case config.Control2:
			if policy.Enabled {
				validateC2(api, policy.Path)
			}
		case config.Control3:
			if policy.Enabled {
				validateC3(api, policy.Path)
			}
		case config.Control4:
			if policy.Enabled {
				validateC4(api, policy.Path)
			}
		}
	}
}

func (h *Handler) ValidateC1(policyPath string) {
	fmt.Println("------------------------------Control-1------------------------------")

	policy := common.UserAuthPolicy(policyPath)
	ciCommits, err := api.Repo.GetChangesToCiCd(
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	if ciCommits != nil {
		verifyCiCdCommitsAuthPolicy(ciCommits, policy, cfg.RepoInfoChecks.TrustedData)
		return
	}
	if err != nil {
		fmt.Printf("[Control 1: ERROR - performing control-1: %v]", err.Error())
		return
	}
	if ciCommits == nil {
		msg := fmt.Sprintf("{ \"control\": \"Control 1\", \"level\": \"WARNING\", \"msg\": \"No new commits since %v\"}", sinceDate)
		fmt.Println(msg)
		common.SendNotification(msg)
	}
}

func (h *Handler) ValidateC2(policyPath string) {
	fmt.Println("------------------------------Control-2------------------------------")

	signatureProtection := api.GetProjectSignatureProtection()
	policy := common.SignatureProtectionPolicy(policyPath)
	verifyRepoProtectionPolicy(&signatureProtection, policy)
}

func (h *Handler) ValidateC3(api *gitlab.Api, policyPath string) {
	fmt.Println("------------------------------Control-3------------------------------")

	automationKeys, _ := api.GetAutomationKeys()

	policy := common.KeyExpiryPolicy(policyPath)
	verifyExpiryKeysPolicy(automationKeys, policy)
}

func (h *Handler) ValidateC4(api *gitlab.Api, policyPath string) {
	fmt.Println("------------------------------Control-4------------------------------")
	automationKeys, _ := api.GetAutomationKeys()

	policy := common.KeyReadOnlyPolicy(policyPath)
	verifyExpiryKeysPolicy(automationKeys, policy)
}

func verifyCiCdCommitsAuthPolicy(commits []gitlab.CommitInfo, policy common.Policy, data map[string]interface{}) {
	pr := common.CreateRegoWithDataStorage(policy, data)

	for _, commit := range commits {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(commit))
		common.SendNotification(evaluation)
	}
}

func verifyRepoProtectionPolicy(repoProtection *gitlab.RepoCommitProtection, policy common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(repoProtection))
	common.SendNotification(evaluation)
}

func verifyExpiryKeysPolicy(automationKeys []gitlab.AutomationKey, policy common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	for _, automationKey := range automationKeys {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(automationKey))
		common.SendNotification(evaluation)
	}
}
