package github

import (
	"fmt"
	"secure-pipeline-poc/app/clients/github"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

func GitHubUserAuthPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/github/c1_github_user_auth.rego",
		Query: "data.github.user.cicd.auth.is_authorized",
	}
}

func GitHubBranchProtectionPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/github/c2_github_branch_protection.rego",
		Query: "data.github.branch.protection.is_protected",
	}
}

func GitHubKeyExpiryPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/github/c3_github_token_expiry.rego",
		Query: "data.github.token.expiry.needs_update",
	}
}

func GitHubKeyReadOnlyPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/github/c4_github_keys_readonly.rego",
		Query: "data.github.keys.readonly.is_read_only",
	}
}

func ValidatePolicies(token string, cfg config.Config, sinceDate time.Time) {
	gitHubClient := github.NewClient(token)
	fmt.Println("------------------------------Control-1------------------------------")

	// Control-1
	var c1Policy = GitHubUserAuthPolicy()

	var trustedData = config.LoadFileToJsonMap(cfg.RepoInfoChecks.TrustedDataFile)
	ciCommits, _ := github.GetChangesToCiCd(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	verifyCiCdCommitsAuthtPolicy(ciCommits, &c1Policy, trustedData)

	fmt.Println("------------------------------Control-2------------------------------")

	// Control-2
	var c2Policy = GitHubBranchProtectionPolicy()
	signatureProtection := github.GetBranchSignatureProtection(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.ProtectedBranches,
	)
	verifyBranchProtectionPolicy(signatureProtection, &c2Policy)

	fmt.Println("------------------------------Control-3------------------------------")

	// Control-3
	var c3Policy = GitHubKeyExpiryPolicy()
	automationKeysE, _ := github.GetAutomationKeysExpiry(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	verifyExpiryKeysPolicy(automationKeysE, &c3Policy)

	fmt.Println("------------------------------Control-4------------------------------")

	// Control-4
	var c4Policy = GitHubKeyReadOnlyPolicy()
	automationKeysRO, _ := github.GetAutomationKeysPermissions(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	verifyReadOnlyKeysPolicy(automationKeysRO, &c4Policy)
}

func verifyCiCdCommitsAuthtPolicy(commits []github.CommitInfo, policy *common.Policy, data map[string]interface{}) {
	pr := common.CreateRegoWithDataStorage(policy, data)

	for _, commit := range commits {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(commit))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyBranchProtectionPolicy(branchesProtection []github.BranchCommitProtection, policy *common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	for _, branchProtection := range branchesProtection {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(branchProtection))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyExpiryKeysPolicy(automationKeys []github.AutomationKey, policy *common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	for _, automationKey := range automationKeys {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(automationKey))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyReadOnlyKeysPolicy(automationKeys []github.AutomationKey, policy *common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	for _, automationKey := range automationKeys {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(automationKey))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

