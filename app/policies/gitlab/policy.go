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

func ValidateC1(api *gitlab.Api, cfg *config.Config, policyPath string, sinceDate time.Time) {
	fmt.Println("------------------------------Control-1------------------------------")

	policy := common.UserAuthPolicy(policyPath)
	ciCommits, _ := api.Repo.GetChangesToCiCd(
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	verifyCiCdCommitsAuthPolicy(ciCommits, policy, cfg.RepoInfoChecks.TrustedData)
}

func validateC2(api *gitlab.Api, policyPath string) {
	fmt.Println("------------------------------Control-2------------------------------")

	signatureProtection := api.GetProjectSignatureProtection()
	policy := common.SignatureProtectionPolicy(policyPath)
	verifyRepoProtectionPolicy(&signatureProtection, policy)
}

func validateC3(api *gitlab.Api, policyPath string) {
	fmt.Println("------------------------------Control-3------------------------------")

	automationKeys, _ := api.GetAutomationKeys()

	policy := common.KeyExpiryPolicy(policyPath)
	verifyExpiryKeysPolicy(automationKeys, policy)
}

func validateC4(api *gitlab.Api, policyPath string) {
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
