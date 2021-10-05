package gitlab

import (
	"fmt"
	"path"
	"runtime"
	"secure-pipeline-poc/app/clients/gitlab"
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
		PolicyFile: fmt.Sprintf("%vc1_gitlab_user_auth.rego", currDir()),
		Query:      "data.gitlab.user.cicd.auth.is_authorized",
	}
}

func RepoProtectionPolicy() common.Policy {
	return common.Policy{
		PolicyFile: fmt.Sprintf("%vc2_gitlab_repo_protection.rego", currDir()),
		Query:      "data.gitlab.repo.protection.is_protected",
	}
}

func keyExpiryPolicy() common.Policy {
	return common.Policy{
		PolicyFile: fmt.Sprintf("%vc3_gitlab_token_expiry.rego", currDir()),
		Query:      "data.gitlab.token.expiry.needs_update",
	}
}

func keyReadOnlyPolicy() common.Policy {
	return common.Policy{
		PolicyFile: fmt.Sprintf("%vc4_gitlab_keys_readonly.rego", currDir()),
		Query:      "data.gitlab.keys.readonly.is_read_only",
	}
}

func ValidatePolicies(token string, cfg *config.Config, sinceDate time.Time) {
	api := gitlab.NewApi(token, cfg)

	for _, control := range cfg.RepoInfoChecks.ControlsToRun {
		switch control {
		case config.Control1:
			ValidateC1(api, cfg, sinceDate)
		case config.Control2:
			validateC2(api)
		case config.Control3:
			validateC3(api)
		case config.Control4:
			validateC4(api)
		}
	}
}

func ValidateC1(api *gitlab.Api, cfg *config.Config, sinceDate time.Time) {
	fmt.Println("------------------------------Control-1------------------------------")

	policy := userAuthPolicy()
	trustedData := common.LoadFileToJsonMap(cfg.RepoInfoChecks.TrustedDataFile)

	ciCommits, _ := api.Repo.GetChangesToCiCd(
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	verifyCiCdCommitsAuthtPolicy(ciCommits, policy, trustedData)
}

func validateC2(api *gitlab.Api) {
	fmt.Println("------------------------------Control-2------------------------------")

	signatureProtection := api.GetProjectSignatureProtection()
	policy := RepoProtectionPolicy()
	verifyRepoProtectionPolicy(&signatureProtection, policy)
}

func validateC3(api *gitlab.Api) {
	fmt.Println("------------------------------Control-3------------------------------")

	automationKeys, _ := api.GetAutomationKeys()

	policy := keyExpiryPolicy()
	verifyExpiryKeysPolicy(automationKeys, policy)
}

func validateC4(api *gitlab.Api) {
	fmt.Println("------------------------------Control-4------------------------------")
	automationKeys, _ := api.GetAutomationKeys()

	policy := keyReadOnlyPolicy()
	verifyExpiryKeysPolicy(automationKeys, policy)
}

func verifyCiCdCommitsAuthtPolicy(commits []gitlab.CommitInfo, policy common.Policy, data map[string]interface{}) {
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
