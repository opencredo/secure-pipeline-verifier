package gitlab

import (
	"fmt"
	x "github.com/xanzy/go-gitlab"
	c "secure-pipeline-poc/app/clients/gitlab"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"time"
)

func userAuthPolicy() *common.Policy {
	return &common.Policy{
		PolicyFile: "app/policies/c1_gitlab_user_auth.rego",
		Query: "data.gitlab.user.cicd.auth.is_authorized",
	}
}

func branchProtectionPolicy() common.Policy {
	return common.Policy{
		PolicyFile: "app/policies/c2_gitlab_branch_protection.rego",
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

	validateC1(client, cfg, sinceDate)

}


func validateC1(client *x.Client, cfg *config.Config, sinceDate time.Time){
	fmt.Println("------------------------------Control-1------------------------------")

	authPolicy := userAuthPolicy()
	trustedData := config.LoadFileToJsonMap(cfg.RepoInfoChecks.TrustedDataFile)

	ciCommits, _ := c.GetChangesToCiCd(
		client,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	common.VerifyCiCdCommitsAuthtPolicy(ciCommits, authPolicy, trustedData)

}