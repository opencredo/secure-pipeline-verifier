package gitlab

import (
	"fmt"
	x "github.com/xanzy/go-gitlab"
	"secure-pipeline-poc/app/clients/gitlab"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

func userAuthPolicy() *common.Policy {
	return &common.Policy{
		PolicyFile: "app/policies/gitlab/c1_gitlab_user_auth.rego",
		Query: "data.gitlab.user.cicd.auth.is_authorized",
	}
}

func branchProtectionPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/gitlab/c2_gitlab_branch_protection.rego",
		Query: "data.gitlab.branch.protection.is_protected",
	}
}

func keyExpiryPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/c3_gitlab_token_expiry.rego",
		Query: "data.gitlab.token.expiry.needs_update",
	}
}

func keyReadOnlyPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/c4_gitlab_keys_readonly.rego",
		Query: "data.gitlab.keys.readonly.is_read_only",
	}
}

func ValidatePolicies(token string, cfg *config.Config, sinceDate time.Time) {
	client, _ := x.NewClient(token)
	// Endpoint to the project
	projectPath := fmt.Sprintf("%s/%s", cfg.Project.Owner, cfg.Project.Repo)

	validateC1(client, cfg, projectPath, sinceDate)
	validateC2(client, cfg, projectPath)
}

func validateC1(client *x.Client, cfg *config.Config, projectPath string, sinceDate time.Time){
	fmt.Println("------------------------------Control-1------------------------------")

	policy := userAuthPolicy()
	trustedData := config.LoadFileToJsonMap(cfg.RepoInfoChecks.TrustedDataFile)

	ciCommits, _ := gitlab.GetChangesToCiCd(
		client,
		projectPath,
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	verifyCiCdCommitsAuthtPolicy(ciCommits, policy, trustedData)
}

func validateC2(client *x.Client, cfg *config.Config, projectPath string){
	fmt.Println("------------------------------Control-2------------------------------")

	signatureProtection := gitlab.GetBranchSignatureProtection(client, projectPath, cfg.RepoInfoChecks.ProtectedBranches)
	policy := branchProtectionPolicy()
	verifyBranchProtectionPolicy(signatureProtection, policy)
}

func verifyCiCdCommitsAuthtPolicy(commits []gitlab.CommitInfo, policy *common.Policy, data map[string]interface{}) {
	pr := common.CreateRegoWithDataStorage(policy, data)

	for _, commit := range commits {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(commit))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyBranchProtectionPolicy(branchesProtection []gitlab.BranchCommitProtection, policy common.Policy) {
	pr := common.CreateRegoWithoutDataStorage(policy)

	for _, branchProtection := range branchesProtection {
		evaluation := common.EvaluatePolicy(pr, common.GetObjectMap(branchProtection))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}