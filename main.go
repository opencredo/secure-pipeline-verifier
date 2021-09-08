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
		fmt.Sprintf("Error %v occurred while parsing date from %v", err.Error(), os.Args[2])
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

	verifyCiCdCommitsAuthtPolicy(ciCommits, c1Policy, trustedData)

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

func verifyCiCdCommitsAuthtPolicy(commits []client.CommitInfo, policy policy.Policy, data map[string]interface{}) {
	pr := createRegoWithDataStorage(policy, data)

	for _, commit := range commits {
		evaluation := evaluatePolicy(pr, getObjectMap(commit))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyBranchProtectionPolicy(branchesProtection []client.BranchCommitProtection, policy policy.Policy) {
	pr := createRegoWithoutDataStorage(policy)

	for _, branchProtection := range branchesProtection {
		evaluation := evaluatePolicy(pr, getObjectMap(branchProtection))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyExpiryKeysPolicy(automationKeys []client.AutomationKey, policy policy.Policy) {
	pr := createRegoWithoutDataStorage(policy)

	for _, automationKey := range automationKeys {
		evaluation := evaluatePolicy(pr, getObjectMap(automationKey))
		// send the info/warning message to Slack
		fmt.Println("", evaluation)
	}
}

func verifyReadOnlyKeysPolicy(automationKeys []client.AutomationKey, policy policy.Policy) {
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

func createRegoWithoutDataStorage(policy policy.Policy) rego.PartialResult {
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

func createRegoWithDataStorage(policy policy.Policy, data map[string]interface{}) rego.PartialResult {
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
