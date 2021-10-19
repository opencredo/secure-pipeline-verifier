package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"os"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/notification"
	"time"
)

type Policy struct {
	PolicyFile string
	Query      string
}

func UserAuthPolicy(path string) *Policy {
	return &Policy{
		PolicyFile: path,
		Query:      "data.user.cicd.auth.is_authorized",
	}
}

func SignatureProtectionPolicy(path string) *Policy {
	return &Policy{
		PolicyFile: path,
		Query:      "data.signature.protection.is_protected",
	}
}

func KeyExpiryPolicy(path string) *Policy {
	return &Policy{
		PolicyFile: path,
		Query:      "data.token.expiry.needs_update",
	}
}

func KeyReadOnlyPolicy(path string) *Policy {
	return &Policy{
		PolicyFile: path,
		Query:      "data.keys.readonly.is_read_only",
	}
}

func (p *Policy) Process(slackCfg config.Slack, input map[string]interface{}, dataStorage ...map[string]interface{}) {

	if dataStorage != nil {
		CreateRegoWithDataStorage(p, dataStorage[0])
	}

	pr := CreateRegoWithoutDataStorage(p)
	var evals []interface{}
	for _, item := range input {
		evaluation := EvaluatePolicy(pr, GetObjectMap(item))

		evals = append(evals, evaluation)

		fmt.Println("", evaluation)
	}
	// send the info/warning message to Slack
	SendNotification(evals, slackCfg)
}

type PoliciesReader interface {
	UserAuthPolicy(path string) *Policy
	SignatureProtectionPolicy(path string) *Policy
	KeyExpiryPolicy(path string) *Policy
	KeyReadOnlyPolicy(path string) *Policy
}

type Controls interface {
	SetClient(token string)
	ValidateC1(policyPath string, cfg *config.Config, sinceDate time.Time)
	ValidateC2(policyPath string, cfg *config.Config)
	ValidateC3(policyPath string, cfg *config.Config)
	ValidateC4(policyPath string, cfg *config.Config)
}

type ValidateInput struct {
	Config   *config.Config
	Controls Controls
	SinceDate time.Time
	Token     string
}

func ValidatePolicies(i *ValidateInput) {
	i.Controls.SetClient(i.Token)

	for _, policy := range i.Config.RepoInfoChecks.Policies {
		switch policy.Control {
		case config.Control1:
			if policy.Enabled {
				i.Controls.ValidateC1(policy.Path, i.Config, i.SinceDate)
			}
		case config.Control2:
			if policy.Enabled {
				i.Controls.ValidateC2(policy.Path, i.Config)
			}
		case config.Control3:
			if policy.Enabled {
				i.Controls.ValidateC3(policy.Path, i.Config)
			}
		case config.Control4:
			if policy.Enabled {
				i.Controls.ValidateC4(policy.Path, i.Config)
			}
		}
	}
}

func CreateRegoWithDataStorage(policy *Policy, data map[string]interface{}) *rego.PartialResult {
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

func CreateRegoWithoutDataStorage(policy *Policy) *rego.PartialResult {
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
	fmt.Printf("Json Response: %s \n", jsonObject)
	var objectMap map[string]interface{}
	_ = json.Unmarshal(jsonObject, &objectMap)
	return objectMap
}

func SendNotification(evaluation interface{}, slackConfig config.Slack) {
	evalMap := evaluation.(map[string]interface{})
	fmt.Println("Evaluation:")
	fmt.Println(" - Control: ", evalMap["control"].(string))
	fmt.Println(" - Level: ", evalMap["level"].(string))
	fmt.Println(" - Message: ", evalMap["msg"].(string))

	notification.Notify(evaluation, slackConfig)
}
