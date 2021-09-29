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

func userAuthPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/gitlab/c1_gitlab_user_auth.rego",
		Query: "data.gitlab.user.cicd.auth.is_authorized",
	}
}

func RepoProtectionPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/gitlab/c2_gitlab_repo_protection.rego",
		Query: "data.gitlab.repo.protection.is_protected",
	}
}

func keyExpiryPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/gitlab/c3_gitlab_token_expiry.rego",
		Query: "data.gitlab.token.expiry.needs_update",
	}
}

func keyReadOnlyPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/gitlab/c4_gitlab_keys_readonly.rego",
		Query: "data.gitlab.keys.readonly.is_read_only",
	}
}

func ValidatePolicies(token string, cfg *config.Config, sinceDate time.Time) {
	client, _ := x.NewClient(token)
	// Endpoint to the project
	projectPath := fmt.Sprintf("%s/%s", cfg.Project.Owner, cfg.Project.Repo)

	validateC1(client, cfg, projectPath, sinceDate)
	validateC2(client, projectPath)

	automationKeys, _ := gitlab.GetAutomationKeys(client, projectPath)

	validateC3(automationKeys)
	validateC4(automationKeys)
}

func validateC1(client *x.Client, cfg *config.Config, projectPath string, sinceDate time.Time){
	fmt.Println("------------------------------Control-1------------------------------")

	policy := userAuthPolicy()
	trustedData := common.LoadFileToJsonMap(cfg.RepoInfoChecks.TrustedDataFile)

	ciCommits, _ := gitlab.GetChangesToCiCd(
		client,
		projectPath,
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	verifyCiCdCommitsAuthtPolicy(ciCommits, policy, trustedData)
}

func validateC2(client *x.Client, projectPath string){
	fmt.Println("------------------------------Control-2------------------------------")

	signatureProtection := gitlab.GetProjectSignatureProtection(client, projectPath)
	policy := RepoProtectionPolicy()
	verifyRepoProtectionPolicy(&signatureProtection, policy)
}

func validateC3(automationKeys []gitlab.AutomationKey){
	fmt.Println("------------------------------Control-3------------------------------")

	policy := keyExpiryPolicy()
	verifyExpiryKeysPolicy(automationKeys, policy)
}

func validateC4(automationKeys []gitlab.AutomationKey){
	fmt.Println("------------------------------Control-4------------------------------")

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