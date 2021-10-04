package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)


const (
	GitHubToken = "GITHUB_TOKEN"	// Env Variable - Token to call GitHub REST APIs
	GitLabToken = "GITLAB_TOKEN"	// Env Variable - Token to call GitLab REST APIs
	SlackToken  = "SLACK_TOKEN"		// Env Variable - Token to connect to Slack for notifications

	Control1 = "c1"
	Control2 = "c2"
	Control3 = "c3"
	Control4 = "c4"
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
	ControlsToRun 	  []string `yaml:"controls-to-run"`
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
