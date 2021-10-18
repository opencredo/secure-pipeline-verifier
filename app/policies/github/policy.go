package github

import (
	"fmt"
	x "github.com/google/go-github/v38/github"
	"secure-pipeline-poc/app/clients/github"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

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

	var policy = common.UserAuthPolicy(policyPath)

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
		msg := fmt.Sprintf("{ \"control\": \"Control 1\", \"level\": \"WARNING\", \"msg\": \"No new commits since %v\"}", sinceDate)
		fmt.Println(msg)
		common.SendNotification(msg)
	}
}

func validateC2(client *x.Client, policyPath string, cfg *config.Config) {
	fmt.Println("------------------------------Control-2------------------------------")

	var c2Policy = common.SignatureProtectionPolicy(policyPath)
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

	var policy = common.KeyExpiryPolicy(policyPath)
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

	var policy = common.KeyReadOnlyPolicy(policyPath)
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

	for _, commit := range commits {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(commit))
		common.SendNotification(evaluation)
	}
}

func verifyBranchProtectionPolicy(branchesProtection []github.BranchCommitProtection, policy common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	for _, branchProtection := range branchesProtection {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(branchProtection))
		common.SendNotification(evaluation)
	}
}

func verifyExpiryKeysPolicy(automationKeys []github.AutomationKey, policy common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	for _, automationKey := range automationKeys {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(automationKey))
		common.SendNotification(evaluation)
	}
}

func verifyReadOnlyKeysPolicy(automationKeys []github.AutomationKey, policy common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	for _, automationKey := range automationKeys {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(automationKey))
		common.SendNotification(evaluation)
	}
}
