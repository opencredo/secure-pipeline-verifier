package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

type Policy struct {
	PolicyFile string
	Query string
}

func GitLabUserAuthPolicy() Policy {
	return Policy{
		PolicyFile: "app/policies/c1_gitlab_user_auth.rego",
		Query: "data.gitlab.user.cicd.auth.is_authorized",
	}
}

func GitLabBranchProtectionPolicy() Policy {
	return Policy{
		PolicyFile: "app/policies/c2_gitlab_branch_protection.rego",
		Query: "data.gitlab.branch.protection.is_protected",
	}
}

func GitLabKeyExpiryPolicy() Policy {
	return Policy{
		PolicyFile: "app/policies/c3_gitlab_token_expiry.rego",
		Query: "data.gitlab.token.expiry.needs_update",
	}
}

func GitLabKeyReadOnlyPolicy() Policy {
	return Policy{
		PolicyFile: "app/policies/c4_gitlab_keys_readonly.rego",
		Query: "data.gitlab.keys.readonly.is_read_only",
	}
}

func ValidatePolicies(token string) {
	client, _ := gitlab.NewClient(token)

	validateC1(client)

}


func validateC1(client *gitlab.Client){
	fmt.Println("------------------------------Control-1------------------------------")
	
}