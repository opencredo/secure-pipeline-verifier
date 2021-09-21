package github

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"io/ioutil"
	"os"
	"secure-pipeline-poc/app/clients/github"
	"secure-pipeline-poc/app/config"
	"time"
)

type Policy struct {
	PolicyFile string
	Query string
}

func GitHubUserAuthPolicy() Policy {
	return Policy{
		PolicyFile: "app/policies/github/c1_github_user_auth.rego",
		Query: "data.github.user.cicd.auth.is_authorized",
	}
}

func GitHubBranchProtectionPolicy() Policy {
	return Policy{
		PolicyFile: "app/policies/github/c2_github_branch_protection.rego",
		Query: "data.github.branch.protection.is_protected",
	}
}

func GitHubKeyExpiryPolicy() Policy {
	return Policy{
		PolicyFile: "app/policies/github/c3_github_token_expiry.rego",
		Query: "data.github.token.expiry.needs_update",
	}
}

func GitHubKeyReadOnlyPolicy() Policy {
	return Policy{
		PolicyFile: "app/policies/github/c4_github_keys_readonly.rego",
		Query: "data.github.keys.readonly.is_read_only",
	}
}

func ValidatePolicies(token string, cfg config.Config, sinceDate time.Time) {
	gitHubClient := github.NewClient(token)
	fmt.Println("------------------------------Control-1------------------------------")

	// Control-1
	var c1Policy = GitHubUserAuthPolicy()

	var trustedData = loadFileToJsonMap(cfg.RepoInfoChecks.TrustedDataFile)
	ciCommits, _ := github.GetChangesToCiCd(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.CiCdPath,
		sinceDate,
	)

	verifyCiCdCommitsAuthtPolicy(ciCommits, c1Policy, trustedData)

	fmt.Println("------------------------------Control-2------------------------------")

	// Control-2
	var c2Policy = GitHubBranchProtectionPolicy()
	signatureProtection := github.GetBranchSignatureProtection(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
		cfg.RepoInfoChecks.ProtectedBranches,
	)
	verifyBranchProtectionPolicy(signatureProtection, c2Policy)

	fmt.Println("------------------------------Control-3------------------------------")

	// Control-3
	var c3Policy = GitHubKeyExpiryPolicy()
	automationKeysE, _ := github.GetAutomationKeysExpiry(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	verifyExpiryKeysPolicy(automationKeysE, c3Policy)

	fmt.Println("------------------------------Control-4------------------------------")

	// Control-4
	var c4Policy = GitHubKeyReadOnlyPolicy()
	automationKeysRO, _ := github.GetAutomationKeysPermissions(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	verifyReadOnlyKeysPolicy(automationKeysRO, c4Policy)
}

func verifyCiCdCommitsAuthtPolicy(commits []github.CommitInfo, policy Policy, data map[string]interface{}) {
	pr := createRegoWithDataStorage(policy, data)

	for _, commit := range commits {
		evaluation := evaluatePolicy(pr, getObjectMap(commit))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyBranchProtectionPolicy(branchesProtection []github.BranchCommitProtection, policy Policy) {
	pr := createRegoWithoutDataStorage(policy)

	for _, branchProtection := range branchesProtection {
		evaluation := evaluatePolicy(pr, getObjectMap(branchProtection))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyExpiryKeysPolicy(automationKeys []github.AutomationKey, policy Policy) {
	pr := createRegoWithoutDataStorage(policy)

	for _, automationKey := range automationKeys {
		evaluation := evaluatePolicy(pr, getObjectMap(automationKey))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyReadOnlyKeysPolicy(automationKeys []github.AutomationKey, policy Policy) {
	pr := createRegoWithoutDataStorage(policy)

	for _, automationKey := range automationKeys {
		evaluation := evaluatePolicy(pr, getObjectMap(automationKey))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func getObjectMap(anObject interface{}) map[string]interface{} {
	jsonObject, _ := json.MarshalIndent(anObject, "", "  ")
	fmt.Printf("Json: %s \n", jsonObject)
	var objectMap map[string]interface{}
	_ = json.Unmarshal(jsonObject, &objectMap)
	return objectMap
}

func createRegoWithoutDataStorage(policy Policy) rego.PartialResult {
	ctx := context.Background()
	r := rego.New(
		rego.Query(policy.Query),
		rego.Load([]string{policy.PolicyFile}, nil),
	)

	pr, err := r.PartialResult(ctx)
	if err != nil {
		fmt.Println("Error occurred while creating partial result. Exiting!", err)
		os.Exit(2)
	}

	return pr
}

func createRegoWithDataStorage(policy Policy, data map[string]interface{}) rego.PartialResult {
	ctx := context.Background()
	store := inmem.NewFromObject(data)

	txn, err := store.NewTransaction(ctx, storage.WriteParams)
	if err != nil {
		panic(err)
	}

	r := rego.New(
		rego.Query(policy.Query),
		rego.Store(store),
		rego.Transaction(txn),
		rego.Load([]string{policy.PolicyFile}, nil),
	)

	pr, err := r.PartialResult(ctx)
	if err != nil {
		fmt.Println("Error occurred while creating partial result. Exiting!", err)
		os.Exit(2)
	}

	return pr
}

func evaluatePolicy(pr rego.PartialResult, commit map[string]interface{}) string {
	ctx := context.Background()

	r := pr.Rego(
		rego.Input(commit),
	)

	// Run evaluation.
	rs, err := r.Eval(ctx)
	if err != nil {
		fmt.Println("Error evaluating policy", err)
	}

	return fmt.Sprintf("%v",rs[0].Expressions[0].Value)

}

func loadFileToJsonMap(filename string) map[string]interface{} {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteContent, _ := ioutil.ReadAll(jsonFile)

	var content map[string]interface{}
	_ = json.Unmarshal(byteContent, &content)

	return content
}
