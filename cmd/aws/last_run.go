// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"time"
)

// SSMGetParameterAPI defines the interface for the GetParameter function.
// We use this interface to test the function using a mocked service.
type SSMGetParameterAPI interface {
	GetParameter(ctx context.Context,
		params *ssm.GetParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

// SSMPutParameterAPI defines the interface for the PutParameter function.
// We use this interface to test the function using a mocked service.
type SSMPutParameterAPI interface {
	PutParameter(ctx context.Context,
		params *ssm.PutParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
}

// FindParameter retrieves an AWS Systems Manager string parameter
func FindParameter(c context.Context, api SSMGetParameterAPI, value *string) (*ssm.GetParameterOutput, error) {
	input := &ssm.GetParameterInput{
		Name: value,
	}
	return api.GetParameter(c, input)
}

// PutParameter updates a key in AWS Systems Manager
func PutParameter(c context.Context, api SSMPutParameterAPI, key *string, value *string) (*ssm.PutParameterOutput, error) {
	input := &ssm.PutParameterInput{
		Name:  key,
		Value: value,
		Type:  types.ParameterTypeString,
	}
	return api.PutParameter(c, input)
}

// ProcessLastRun saves last run in Parameter Store
func ProcessLastRun(ctx context.Context, cfg aws.Config) (string, error){
	client := ssm.NewFromConfig(cfg)

	parameterKey := "last_run"

	getParam, err := FindParameter(ctx, client, &parameterKey)
	// The endpoint returns error when the parameter doesn't exist
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	lastRun := *getParam.Parameter.Value
	now := time.Now().String()

	if lastRun == " " {
		fmt.Println("No previous runs found.")
	} else {
		fmt.Printf("Last run was at %v", lastRun)
	}

	_, err = PutParameter(ctx, client, &parameterKey, &now)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return "", nil

}
