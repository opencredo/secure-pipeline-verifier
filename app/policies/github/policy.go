package github

import (
	"fmt"
	x "github.com/google/go-github/v38/github"
	"secure-pipeline-poc/app/clients/github"
	"secure-pipeline-poc/app/notification"
	"secure-pipeline-poc/app/policies/common"
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

type Handler struct {
	Client *x.Client
	Platform *common.Platform
}

func (h *Handler) SetClient(token string) {
	h.Client = github.NewClient(token)
}

func (h *Handler) ValidateC1(policyPath string) {
	fmt.Println("------------------------------Control-1------------------------------")

	var policy = userAuthPolicy(policyPath)

	ciCommits, err := github.GetChangesToCiCd(
		h.Client,
		h.Platform.Config,
		h.Platform.SinceDate,
	)

	if ciCommits != nil {
		verifyCiCdCommitsAuthPolicy(ciCommits, policy, h.Platform.Config.RepoInfoChecks.TrustedData)
		return
	}
	if err != nil {
		fmt.Printf("[Control 1: ERROR - performing control-1: %v]", err.Error())
		return
	}
	if ciCommits == nil {
		msg := fmt.Sprintf("[Control 1: WARNING - No new commits since %v]", h.Platform.SinceDate)
		fmt.Println(msg)
		notification.Notify([]string{msg})
	}
}

func (h *Handler) ValidateC2(policyPath string) {
	fmt.Println("------------------------------Control-2------------------------------")

	var c2Policy = branchProtectionPolicy(policyPath)
	signatureProtection := github.GetBranchSignatureProtection(
		h.Client,
		h.Platform.Config.Project.Owner,
		h.Platform.Config.Project.Repo,
		h.Platform.Config.RepoInfoChecks.ProtectedBranches,
	)
	verifyBranchProtectionPolicy(signatureProtection, c2Policy)
}

func (h *Handler) ValidateC3(policyPath string) {
	fmt.Println("------------------------------Control-3------------------------------")

	var policy = keyExpiryPolicy(policyPath)
	automationKeysE, err := github.GetAutomationKeysExpiry(
		h.Client,
		h.Platform.Config.Project.Owner,
		h.Platform.Config.Project.Repo,
	)
	if automationKeysE != nil {
		verifyExpiryKeysPolicy(automationKeysE, policy)
	}
	if err != nil {
		fmt.Println("Error performing control-3: ", err)
	}
}

func (h *Handler) ValidateC4(policyPath string) {
	fmt.Println("------------------------------Control-4------------------------------")

	var policy = keyReadOnlyPolicy(policyPath)
	automationKeysRO, err := github.GetAutomationKeysPermissions(
		h.Client,
		h.Platform.Config.Project.Owner,
		h.Platform.Config.Project.Repo,
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
