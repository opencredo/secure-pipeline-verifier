package config_test

import (
	"github.com/stretchr/testify/assert"
	"secure-pipeline-poc/app/config"
	"testing"
)

func TestConfigFileLoaded(t *testing.T) {
	assert := assert.New(t)

	var cfg config.Config
	config.LoadConfig("config_test.yaml", &cfg)

	project := cfg.Project
	assert.Equalf("oc", project.Owner, "project Owner should be oc")
	assert.Equalf("my-app-repo", project.Repo, "project Repo should be my-app-repo")

	repoInfoChecks := cfg.RepoInfoChecks
	assert.Equal("oc-trusted-config.json", repoInfoChecks.TrustedDataFile)
	assert.Equal(".github/workflows", repoInfoChecks.CiCdPath)
	assert.Equal([]string{"master", "develop"}, repoInfoChecks.ProtectedBranches, "they should have the same elements")
}