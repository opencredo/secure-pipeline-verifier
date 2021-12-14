package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigFileLoaded(t *testing.T) {
	assert := assert.New(t)

	var configPath = "./test_data/"
	var cfg Config
	LoadConfig(configPath, &cfg)
	LoadTrustedDataToMap(configPath, &cfg)

	project := cfg.Project
	assert.Equalf("github", project.Platform, "project platform should be github")
	assert.Equalf("oc", project.Owner, "project Owner should be oc")
	assert.Equalf("my-app-repo", project.Repo, "project Repo should be my-app-repo")

	repoInfoChecks := cfg.RepoInfoChecks

	assert.Equal(4, len(repoInfoChecks.Policies))
	policy1 := repoInfoChecks.Policies[0]
	assert.Equal("c1", policy1.Control)
	assert.Equal(true, policy1.Enabled)
	assert.Equal("<path to policies>/auth.rego", policy1.Path)
	policy2 := repoInfoChecks.Policies[1]
	assert.Equal("c2", policy2.Control)
	assert.Equal(true, policy2.Enabled)
	assert.Equal("<path to policies>/signed-commits.rego", policy2.Path)
	policy3 := repoInfoChecks.Policies[2]
	assert.Equal("c3", policy3.Control)
	assert.Equal(true, policy3.Enabled)
	assert.Equal("<path to policies>/auth-key-expiry.rego", policy3.Path)
	policy4 := repoInfoChecks.Policies[3]
	assert.Equal("c4", policy4.Control)
	assert.Equal(true, policy4.Enabled)
	assert.Equal("<path to policies>/auth-key-read-only.rego", policy4.Path)

	assert.NotNil(cfg.RepoInfoChecks.TrustedData)
	assert.Equal("some-org/some-repo", cfg.RepoInfoChecks.TrustedData["config"].(map[string]interface{})["repo"])
	assert.Equal("travis", cfg.RepoInfoChecks.TrustedData["config"].(map[string]interface{})["pipeline_type"])
	assert.Equal([]interface{}{"somebody", "someone-else"}, cfg.RepoInfoChecks.TrustedData["config"].(map[string]interface{})["trusted_users"])

	assert.Equal(".github/workflows", repoInfoChecks.CiCdPath)
	assert.Equal([]string{"master", "develop"}, repoInfoChecks.ProtectedBranches, "they should have the same elements")

	assert.Equal(true, cfg.Notifications.Slack.Enabled)
	assert.Equal("INFO", cfg.Notifications.Slack.Level)
	assert.Equal("secure-pipeline", cfg.Notifications.Slack.Channel)
}
