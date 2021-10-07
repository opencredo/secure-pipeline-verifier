package gitlab

import (
	"fmt"
	"secure-pipeline-poc/app/clients/gitlab"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/notification"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

func userAuthPolicy(path string) common.Policy {
	return common.Policy{
		PolicyFile: path,
		Query:      "data.gitlab.user.cicd.auth.is_authorized",
	}
}

func RepoProtectionPolicy(path string) common.Policy {
	return common.Policy{
		PolicyFile: path,
		Query:      "data.gitlab.repo.protection.is_protected",
	}
}

func keyExpiryPolicy(path string) common.Policy {
	return common.Policy{
		PolicyFile: path,
		Query:      "data.gitlab.token.expiry.needs_update",
	}
}

func keyReadOnlyPolicy(path string) common.Policy {
	return common.Policy{
		PolicyFile: path,
		Query:      "data.gitlab.keys.readonly.is_read_only",
	}
}

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

	policy := userAuthPolicy(policyPath)
	ciCommits, _ := api.Repo.GetChangesToCiCd(
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	verifyCiCdCommitsAuthPolicy(ciCommits, policy, cfg.RepoInfoChecks.TrustedData)
}

func validateC2(api *gitlab.Api, policyPath string) {
	fmt.Println("------------------------------Control-2------------------------------")

	signatureProtection := api.GetProjectSignatureProtection()
	policy := RepoProtectionPolicy(policyPath)
	verifyRepoProtectionPolicy(&signatureProtection, policy)
}

func validateC3(api *gitlab.Api, policyPath string) {
	fmt.Println("------------------------------Control-3------------------------------")

	automationKeys, _ := api.GetAutomationKeys()

	policy := keyExpiryPolicy(policyPath)
	verifyExpiryKeysPolicy(automationKeys, policy)
}

func validateC4(api *gitlab.Api, policyPath string) {
	fmt.Println("------------------------------Control-4------------------------------")
	automationKeys, _ := api.GetAutomationKeys()

	policy := keyReadOnlyPolicy(policyPath)
	verifyExpiryKeysPolicy(automationKeys, policy)
}

func verifyCiCdCommitsAuthPolicy(commits []gitlab.CommitInfo, policy common.Policy, data map[string]interface{}) {
	pr := common.CreateRegoWithDataStorage(policy, data)
	var messages []string
	for _, commit := range commits {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(commit))

		messages = append(messages, evaluation)
		fmt.Println("", evaluation)
	}
	// send the info/warning message to Slack
	notification.Notify(messages)
}

func verifyRepoProtectionPolicy(repoProtection *gitlab.RepoCommitProtection, policy common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(repoProtection))

	var messages []string
	messages = append(messages, evaluation)
	// send the info/warning message to Slack
	notification.Notify(messages)

	fmt.Println("", evaluation)
}

func verifyExpiryKeysPolicy(automationKeys []gitlab.AutomationKey, policy common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)
	var messages []string
	for _, automationKey := range automationKeys {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(automationKey))

		messages = append(messages, evaluation)

		fmt.Println("", evaluation)
	}
	// send the info/warning message to Slack
	notification.Notify(messages)
}
