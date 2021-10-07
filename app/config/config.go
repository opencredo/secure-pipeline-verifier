package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
)

const (
	ConfigsFileName     = "config.yaml"
	TrustedDataFileName = "trusted-data.json"

	GitHubToken = "GITHUB_TOKEN" // Env Variable - Token to call GitHub REST APIs
	GitLabToken = "GITLAB_TOKEN" // Env Variable - Token to call GitLab REST APIs
	SlackToken  = "SLACK_TOKEN"  // Env Variable - Token to connect to Slack for notifications

	Control1 = "c1"
	Control2 = "c2"
	Control3 = "c3"
	Control4 = "c4"
)

type Project struct {
	Platform string `yaml:"platform"`
	Owner    string `yaml:"owner"`
	Repo     string `yaml:"repo"`
}

type RepoInfoChecks struct {
	TrustedData       map[string]interface{}
	CiCdPath          string   `yaml:"ci-cd-path"`
	ProtectedBranches []string `yaml:"protected-branches"`
	ControlsToRun     []string `yaml:"controls-to-run"`
}

type Config struct {
	Project        Project        `yaml:"project"`
	RepoInfoChecks RepoInfoChecks `yaml:"repo-info-checks"`
}

func LoadConfig(filePath string, cfg *Config) {
	file, err := os.Open(filePath + ConfigsFileName)
	if err != nil {
		processError(err)
	}
	defer file.Close()

	DecodeConfigToStruct(file, cfg)
}

func DecodeConfigToStruct(reader io.Reader, cfg *Config) {
	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func LoadTrustedDataToJsonMap(filePath string, cfg *Config) {
	jsonFile, err := os.Open(filePath + TrustedDataFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	DecodeTrustedDataToMap(jsonFile, cfg)
}

func DecodeTrustedDataToMap(reader io.Reader, cfg *Config) {
	byteContent, _ := ioutil.ReadAll(reader)

	var content map[string]interface{}
	_ = json.Unmarshal(byteContent, &content)
	cfg.RepoInfoChecks.TrustedData = content
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
