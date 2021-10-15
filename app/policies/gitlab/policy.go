package gitlab

import (
	"fmt"
	x "github.com/xanzy/go-gitlab"
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

type PolicyCheck interface {
	Verify()
}

type Handler struct {
	Client *x.Client
	Cfg *config.Config
	SinceDate time.Time
}

func (h *Handler) SetClient(token string) {
	panic("implement me")
}

func (h *Handler) ValidateC1(policyPath string) {
	fmt.Println("------------------------------Control-1------------------------------")

	policy := userAuthPolicy(policyPath)
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
		msg := fmt.Sprintf("[Control 1: WARNING - No new commits since %v]", sinceDate)
		fmt.Println(msg)
		notification.Notify([]string{msg})
	}
}

func (h *Handler) ValidateC2(policyPath string) {
	fmt.Println("------------------------------Control-2------------------------------")

	signatureProtection := api.GetProjectSignatureProtection()
	policy := RepoProtectionPolicy(policyPath)
	verifyRepoProtectionPolicy(&signatureProtection, policy)
}

func (h *Handler) ValidateC3(api *gitlab.Api, policyPath string) {
	fmt.Println("------------------------------Control-3------------------------------")

	automationKeys, _ := api.GetAutomationKeys()

	policy := keyExpiryPolicy(policyPath)
	verifyExpiryKeysPolicy(automationKeys, policy)
}

func (h *Handler) ValidateC4(api *gitlab.Api, policyPath string) {
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
