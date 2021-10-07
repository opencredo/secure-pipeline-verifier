package github

import (
	"fmt"
	x "github.com/google/go-github/v38/github"
	"path"
	"runtime"
	"secure-pipeline-poc/app/clients/github"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/notification"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

// currDir returns current directory of the file
func currDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	p := path.Dir(filename) + "/"
	return p
}

func userAuthPolicy() common.Policy {
	return common.Policy{
		PolicyFile: fmt.Sprintf("%vc1_github_user_auth.rego", currDir()),
		Query:      "data.github.user.cicd.auth.is_authorized",
	}
}

func branchProtectionPolicy() common.Policy {
	return common.Policy{
		PolicyFile: fmt.Sprintf("%vc2_github_branch_protection.rego", currDir()),
		Query:      "data.github.branch.protection.is_protected",
	}
}

func keyExpiryPolicy() common.Policy {
	return common.Policy{
		PolicyFile: fmt.Sprintf("%vc3_github_token_expiry.rego", currDir()),
		Query:      "data.github.token.expiry.needs_update",
	}
}

func keyReadOnlyPolicy() common.Policy {
	return common.Policy{
		PolicyFile: fmt.Sprintf("%vc4_github_keys_readonly.rego", currDir()),
		Query:      "data.github.keys.readonly.is_read_only",
	}
}

func ValidatePolicies(token string, cfg *config.Config, sinceDate time.Time) {
	client := github.NewClient(token)

	for _, control := range cfg.RepoInfoChecks.ControlsToRun {
		switch control {
		case config.Control1:
			validateC1(client, cfg, sinceDate)
		case config.Control2:
			validateC2(client, cfg)
		case config.Control3:
			validateC3(client, cfg)
		case config.Control4:
			validateC4(client, cfg)
		}
	}
}

func validateC1(client *x.Client, cfg *config.Config, sinceDate time.Time) {
	fmt.Println("------------------------------Control-1------------------------------")

	var policy = userAuthPolicy()

	ciCommits, errC1 := github.GetChangesToCiCd(
		client,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	if ciCommits != nil {
		verifyCiCdCommitsAuthPolicy(ciCommits, policy, cfg.RepoInfoChecks.TrustedData)
	}
	if errC1 != nil {
		fmt.Println("Error performing control-1: ", errC1.Error())
	}
}

func validateC2(client *x.Client, cfg *config.Config) {
	fmt.Println("------------------------------Control-2------------------------------")

	var c2Policy = branchProtectionPolicy()
	signatureProtection := github.GetBranchSignatureProtection(
		client,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.ProtectedBranches,
	)
	verifyBranchProtectionPolicy(signatureProtection, c2Policy)
}

func validateC3(client *x.Client, cfg *config.Config) {
	fmt.Println("------------------------------Control-3------------------------------")

	var policy = keyExpiryPolicy()
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

func validateC4(client *x.Client, cfg *config.Config) {
	fmt.Println("------------------------------Control-4------------------------------")

	var policy = keyReadOnlyPolicy()
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
