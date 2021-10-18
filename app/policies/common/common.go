package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"os"
	"secure-pipeline-poc/app/notification"
)

type Policy struct {
	PolicyFile string
	Query      string
}

func UserAuthPolicy(path string) Policy {
	return Policy{
		PolicyFile: path,
		Query:      "data.user.cicd.auth.is_authorized",
	}
}

func SignatureProtectionPolicy(path string) Policy {
	return Policy{
		PolicyFile: path,
		Query:      "data.signature.protection.is_protected",
	}
}

func KeyExpiryPolicy(path string) Policy {
	return Policy{
		PolicyFile: path,
		Query:      "data.token.expiry.needs_update",
	}
}

func KeyReadOnlyPolicy(path string) Policy {
	return Policy{
		PolicyFile: path,
		Query:      "data.keys.readonly.is_read_only",
	}
}

func CreateRegoWithDataStorage(policy Policy, data map[string]interface{}) *rego.PartialResult {
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

	return &pr
}

func CreateRegoWithoutDataStorage(policy Policy) *rego.PartialResult {
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

	return &pr
}

func EvaluatePolicy(pr *rego.PartialResult, input map[string]interface{}) interface{} {
	ctx := context.Background()

	r := pr.Rego(
		rego.Input(input),
	)

	// Run evaluation.
	rs, err := r.Eval(ctx)
	if err != nil {
		fmt.Println("Error evaluating policy", err)
	}

	return rs[0].Expressions[0].Value
}

func GetObjectMap(anObject interface{}) map[string]interface{} {
	jsonObject, _ := json.MarshalIndent(anObject, "", "  ")
	fmt.Printf("Json: %s \n", jsonObject)
	var objectMap map[string]interface{}
	_ = json.Unmarshal(jsonObject, &objectMap)
	return objectMap
}

func SendNotification(evaluation interface{}) {
	notification.Notify(evaluation)
}