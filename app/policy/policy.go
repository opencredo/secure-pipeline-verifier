package policy

type Policy struct {
	PolicyFile string
	Query string
}

func GetGitHubUserAuthzPolicy() Policy {
	return Policy{
		PolicyFile: "app/policy/c1_github_user_authz.rego",
		Query: "data.github.cicd.user.authz.allow",
	}
}

func GitHubBranchProtectionPolicy() Policy {
	return Policy{
		PolicyFile: "app/policy/c2_github_branch_protection.rego",
		Query: "data.github.branch.protection.allow",
	}
}

func GitHubTokenExpiryPolicy() Policy {
	return Policy{
		PolicyFile: "app/policy/c3_github_token_expiry.rego",
		Query: "data.github.token.expiry.allow",
	}
}

func GitHubDeployKeysReadOnlyPolicy() Policy {
	return Policy{
		PolicyFile: "app/policy/c4_github_keys_readonly.rego",
		Query: "data.github.keys.readonly.allow",
	}
}