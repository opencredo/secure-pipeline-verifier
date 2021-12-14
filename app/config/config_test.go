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

	policies := cfg.Policies
	assert.Equal(4, len(policies))

	assert.Equal("c1", policies[0].Control)
	assert.Equal(true, policies[0].Enabled)
	assert.Equal("<path to policies>/auth.rego", policies[0].Path)

	assert.Equal("c2", policies[1].Control)
	assert.Equal(true, policies[1].Enabled)
	assert.Equal("<path to policies>/signed-commits.rego", policies[1].Path)

	assert.Equal("c3", policies[2].Control)
	assert.Equal(true, policies[2].Enabled)
	assert.Equal("<path to policies>/auth-key-expiry.rego", policies[2].Path)

	assert.Equal("c4", policies[3].Control)
	assert.Equal(true, policies[3].Enabled)
	assert.Equal("<path to policies>/auth-key-read-only.rego", policies[3].Path)

	repoInfo := cfg.RepoInfo
	assert.NotNil(repoInfo.TrustedData)
	assert.Equal("some-org/some-repo", cfg.RepoInfo.TrustedData["config"].(map[string]interface{})["repo"])
	assert.Equal("travis", cfg.RepoInfo.TrustedData["config"].(map[string]interface{})["pipeline_type"])
	assert.Equal([]interface{}{"somebody", "someone-else"}, cfg.RepoInfo.TrustedData["config"].(map[string]interface{})["trusted_users"])

	assert.Equal(".github/workflows", repoInfo.CiCdPath)
	assert.Equal([]string{"master", "develop"}, repoInfo.ProtectedBranches, "they should have the same elements")

	assert.Equal(true, cfg.Notifications.Slack.Enabled)
	assert.Equal("INFO", cfg.Notifications.Slack.Level)
	assert.Equal("secure-pipeline", cfg.Notifications.Slack.Channel)
}
