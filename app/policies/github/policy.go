package github

import (
	"fmt"
	x "github.com/google/go-github/v38/github"
	"secure-pipeline-poc/app/clients/github"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/notification"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

func userAuthPolicy(path string) common.Policy {
	return common.Policy{
		PolicyFile: path,
		Query:      "data.github.user.cicd.auth.is_authorized",
	}
}

func branchProtectionPolicy(path string) common.Policy {
	return common.Policy{
		PolicyFile: path,
		Query:      "data.github.branch.protection.is_protected",
	}
}

func keyExpiryPolicy(path string) common.Policy {
	return common.Policy{
		PolicyFile: path,
		Query:      "data.github.token.expiry.needs_update",
	}
}

func keyReadOnlyPolicy(path string) common.Policy {
	return common.Policy{
		PolicyFile: path,
		Query:      "data.github.keys.readonly.is_read_only",
	}
}

func ValidatePolicies(token string, cfg *config.Config, sinceDate time.Time) {
	client := github.NewClient(token)

	for _, policy := range cfg.RepoInfoChecks.Policies {
		switch policy.Control {
		case config.Control1:
			if policy.Enabled {
				validateC1(client, cfg, policy.Path, sinceDate)
			}
		case config.Control2:
			if policy.Enabled {
				validateC2(client, policy.Path, cfg)
			}
		case config.Control3:
			if policy.Enabled {
				validateC3(client, policy.Path, cfg)
			}
		case config.Control4:
			if policy.Enabled {
				validateC4(client, policy.Path, cfg)
			}
		}
	}
}

func validateC1(client *x.Client, cfg *config.Config, policyPath string, sinceDate time.Time) {
	fmt.Println("------------------------------Control-1------------------------------")

	var policy = userAuthPolicy(policyPath)

	ciCommits, err := github.GetChangesToCiCd(
		client,
		cfg.Project.Owner,
		cfg.Project.Repo,
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
		msg := fmt.Sprintf("[Control 1: WARNING - No new commits since %v]", sinceDate)
		fmt.Println(msg)
		notification.Notify([]string{msg})
	}
}

func validateC2(client *x.Client, policyPath string, cfg *config.Config) {
	fmt.Println("------------------------------Control-2------------------------------")

	var c2Policy = branchProtectionPolicy(policyPath)
	signatureProtection := github.GetBranchSignatureProtection(
		client,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.ProtectedBranches,
	)
	verifyBranchProtectionPolicy(signatureProtection, c2Policy)
}

func validateC3(client *x.Client, policyPath string, cfg *config.Config) {
	fmt.Println("------------------------------Control-3------------------------------")

	var policy = keyExpiryPolicy(policyPath)
	automationKeysE, err := github.GetAutomationKeysExpiry(
		client,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	if automationKeysE != nil {
		verifyExpiryKeysPolicy(automationKeysE, policy)
	}
	if err != nil {
		fmt.Println("Error performing control-3: ", err)
	}
}

func validateC4(client *x.Client, policyPath string, cfg *config.Config) {
	fmt.Println("------------------------------Control-4------------------------------")

	var policy = keyReadOnlyPolicy(policyPath)
	automationKeysRO, err := github.GetAutomationKeysPermissions(
		client,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	if automationKeysRO != nil {
		verifyReadOnlyKeysPolicy(automationKeysRO, policy)
	}
	if err != nil {
		fmt.Println("Error performing control-4: ", err)
	}
}

func verifyCiCdCommitsAuthPolicy(commits []github.CommitInfo, policy common.Policy, data map[string]interface{}) {
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

func verifyBranchProtectionPolicy(branchesProtection []github.BranchCommitProtection, policy common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)
	var messages []string

	for _, branchProtection := range branchesProtection {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(branchProtection))
		messages = append(messages, evaluation)
		fmt.Println("", evaluation)
	}
	// send the info/warning message to Slack
	notification.Notify(messages)
}

func verifyExpiryKeysPolicy(automationKeys []github.AutomationKey, policy common.Policy) {
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

func verifyReadOnlyKeysPolicy(automationKeys []github.AutomationKey, policy common.Policy) {
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
