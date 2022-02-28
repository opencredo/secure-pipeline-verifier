package main

import (
	"github.com/secure-pipeline-verifier/app/config"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPoliciesPathUpdate(t *testing.T) {
	assert := assert.New(t)

	var configPath = "../../app/config/test_data/"
	var cfg config.Config
	config.LoadConfig(configPath, &cfg)

	updatePoliciesPath(cfg.Policies)
	assert.Equal(true, strings.HasPrefix(cfg.Policies[0].Path, LambdaPoliciesFolder))
	assert.Equal(true, strings.HasPrefix(cfg.Policies[1].Path, LambdaPoliciesFolder))
	assert.Equal(true, strings.HasPrefix(cfg.Policies[2].Path, LambdaPoliciesFolder))
	assert.Equal(true, strings.HasPrefix(cfg.Policies[3].Path, LambdaPoliciesFolder))
}
