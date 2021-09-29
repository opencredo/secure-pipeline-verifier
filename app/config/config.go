package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// Environment variables
const (
	GitHubToken = "GITHUB_TOKEN"
	GitLabToken = "GITLAB_TOKEN"
	SlackToken  = "SLACK_TOKEN"
)

type Project struct {
	Platform string `yaml:"platform"`
	Owner string 	`yaml:"owner"`
	Repo  string 	`yaml:"repo"`
}

type RepoInfoChecks struct {
	TrustedDataFile   string   `yaml:"trusted-data-file"`
	CiCdPath          string   `yaml:"ci-cd-path"`
	ProtectedBranches []string `yaml:"protected-branches"`
}

type Config struct {
	Project        Project        `yaml:"project"`
	RepoInfoChecks RepoInfoChecks `yaml:"repo-info-checks"`
}

func LoadConfig(filename string, cfg *Config) {
	file, err := os.Open(filename)
	if err != nil {
		processError(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
