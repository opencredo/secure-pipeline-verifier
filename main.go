package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"io/ioutil"
	"os"
	"secure-pipeline-poc/app/client"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policy"
	"time"
)

const GitHubToken = "GITHUB_TOKEN"

var token = os.Getenv(GitHubToken)

func main()  {

	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "path-to-config.yaml", "YYYY-MM-ddTHH:mm:ss.SSSZ")
		return
	}

	var cfg config.Config
	config.LoadConfig(os.Args[1], &cfg)

	sinceDate, err := time.Parse(time.RFC3339, os.Args[2])
	if err != nil {
		fmt.Println("Error occurred while parsing date from ", err)
		os.Exit(2)
	}

	gitHubClient := client.NewClient(token)
	fmt.Println("------------------------------Control-1------------------------------")

	// Control-1
	var c1Policy = policy.GitHubUserAuthPolicy()

	var trustedData = loadFileToJsonMap(cfg.RepoInfoChecks.TrustedDataFile)
	ciCommits, err := client.GetChangesToCiCd(
			gitHubClient,
			cfg.Project.Owner,
			cfg.Project.Repo,
			cfg.RepoInfoChecks.CiCdPath,
			sinceDate,
	)

	verifyCommitsAgainstPolicy(ciCommits, c1Policy, trustedData)

	fmt.Println("------------------------------Control-2------------------------------")

	// Control-2
	var c2Policy = policy.GitHubBranchProtectionPolicy()
	signatureProtection := client.GetBranchSignatureProtection(
			gitHubClient,
			cfg.Project.Owner,
			cfg.Project.Repo,
			cfg.RepoInfoChecks.ProtectedBranches,
	)
	verifyBranchProtectionPolicy(signatureProtection, c2Policy)

	fmt.Println("------------------------------Control-3------------------------------")

	// Control-3
	var c3Policy = policy.GitHubKeyExpiryPolicy()
	automationKeysE, err := client.GetAutomationKeysExpiry(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	verifyExpiryKeysPolicy(automationKeysE, c3Policy)

	fmt.Println("------------------------------Control-4------------------------------")

	// Control-4
	var c4Policy = policy.GitHubKeyReadOnlyPolicy()
	automationKeysRO, err := client.GetAutomationKeysPermissions(
		gitHubClient,
		cfg.Project.Owner,
		cfg.Project.Repo,
	)
	verifyReadOnlyKeysPolicy(automationKeysRO, c4Policy)
}

func verifyCommitsAgainstPolicy(commits []client.CommitInfo, policy policy.Policy, data map[string]interface{}) {
	ctx := context.Background()
	r := createRegoWithConfigData(policy, data)
	pr, err := r.PartialResult(ctx)
	if err != nil {
		fmt.Println("Error occurred while creating partial result", err)
		return
	}

	for _, commit := range commits {
		jsonCommit, _ := json.MarshalIndent(commit, "", "  ")

		fmt.Printf("Commit: %s \n", jsonCommit)
		var commitMap map[string]interface{}
		_ = json.Unmarshal(jsonCommit, &commitMap)
		evaluation := evaluatePolicy(pr, commitMap)
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyBranchProtectionPolicy(branchesProtection []client.BranchCommitProtection, policy policy.Policy) {
	ctx := context.Background()
	r := createRegoWithNoConfigData(policy)
	pr, err := r.PartialResult(ctx)
	if err != nil {
		fmt.Println("Error occurred while creating partial result", err)
		return
	}

	for _, branchProtection := range branchesProtection {
		jsonBranchProtection, _ := json.MarshalIndent(branchProtection, "", "  ")

		fmt.Printf("BranchProtection: %s \n", jsonBranchProtection)
		var branchProtectionMap map[string]interface{}
		_ = json.Unmarshal(jsonBranchProtection, &branchProtectionMap)
		evaluation := evaluatePolicy(pr, branchProtectionMap)
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyExpiryKeysPolicy(automationKeys []client.AutomationKey, policy policy.Policy) {
	ctx := context.Background()
	r := createRegoWithNoConfigData(policy)
	pr, err := r.PartialResult(ctx)
	if err != nil {
		fmt.Println("Error occurred while creating partial result", err)
		return
	}

	for _, automationKey := range automationKeys {
		jsonAutomationKey, _ := json.MarshalIndent(automationKey, "", "  ")

		fmt.Printf("Automation Key: %s \n", jsonAutomationKey)
		var automationKeyMap map[string]interface{}
		_ = json.Unmarshal(jsonAutomationKey, &automationKeyMap)
		evaluation := evaluatePolicy(pr, automationKeyMap)
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyReadOnlyKeysPolicy(automationKeys []client.AutomationKey, policy policy.Policy) {
	ctx := context.Background()
	r := createRegoWithNoConfigData(policy)
	pr, err := r.PartialResult(ctx)
	if err != nil {
		fmt.Println("Error occurred while creating partial result", err)
		return
	}

	for _, automationKey := range automationKeys {
		jsonAutomationKey, _ := json.MarshalIndent(automationKey, "", "  ")

		fmt.Printf("Automation Key: %s \n", jsonAutomationKey)
		var automationKeyMap map[string]interface{}
		_ = json.Unmarshal(jsonAutomationKey, &automationKeyMap)
		evaluation := evaluatePolicy(pr, automationKeyMap)
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func createRegoWithNoConfigData(policy policy.Policy) *rego.Rego {
	return rego.New(
		rego.Query(policy.Query),
		rego.Load([]string{policy.PolicyFile}, nil),
	)
}

func createRegoWithConfigData(policy policy.Policy, data map[string]interface{}) *rego.Rego {
	ctx := context.Background()
	store := inmem.NewFromObject(data)

	txn, err := store.NewTransaction(ctx, storage.WriteParams)
	if err != nil {
		panic(err)
	}

	return rego.New(
		rego.Query(policy.Query),
		rego.Store(store),
		rego.Transaction(txn),
		rego.Load([]string{policy.PolicyFile}, nil),
	)
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
		fmt.Print(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteContent, _ := ioutil.ReadAll(jsonFile)

	var content map[string]interface{}
	_ = json.Unmarshal(byteContent, &content)

	return content
}
