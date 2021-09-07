package policy

type Policy struct {
	PolicyFile string
	Query string
}

func GitHubUserAuthPolicy() Policy {
	return Policy{
		PolicyFile: "app/policy/c1_github_user_unauthz.rego",
		Query: "data.github.user.cicd.auth.is_unauthorized",
	}
}

func GitHubBranchProtectionPolicy() Policy {
	return Policy{
		PolicyFile: "app/policy/c2_github_branch_protection.rego",
		Query: "data.github.branch.protection.is_protected",
	}
}

func GitHubKeyExpiryPolicy() Policy {
	return Policy{
		PolicyFile: "app/policy/c3_github_token_expiry.rego",
		Query: "data.github.token.expiry.needs_update",
	}
}

func GitHubKeyReadOnlyPolicy() Policy {
	return Policy{
		PolicyFile: "app/policy/c4_github_keys_readonly.rego",
		Query: "data.github.keys.readonly.can_write",
	}
}