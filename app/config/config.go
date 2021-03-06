package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
)

const (
	GitHubPlatform = "github"
	GitLabPlatform = "gitlab"

	ConfigsFileName     = "config.yaml"
	TrustedDataFileName = "trusted-data.yaml"

	RepoToken  = "REPO_TOKEN"  // Env Variable - Token to call a version control REST APIs
	SlackToken = "SLACK_TOKEN" // Env Variable - Token to connect to Slack for notifications

	Control1 = "c1"
	Control2 = "c2"
	Control3 = "c3"
	Control4 = "c4"
)

var NotificationLevel = map[string]int{
	"INFO":    0,
	"WARNING": 1,
	"ERROR":   2,
}

type Project struct {
	Platform string `yaml:"platform"`
	Owner    string `yaml:"owner"`
	Repo     string `yaml:"repo"`
}

type Policies struct {
	Control string `yaml:"control"`
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type RepoInfo struct {
	TrustedData       map[string]interface{}
	CiCdPath          string   `yaml:"ci-cd-path"`
	ProtectedBranches []string `yaml:"protected-branches"`
}

type Notifications struct {
	Slack Slack `yaml:"slack"`
}

type Slack struct {
	Enabled bool   `yaml:"enabled"`
	Level   string `yaml:"level"`
	Channel string `yaml:"notification-channel"`
}

type Config struct {
	Project       Project       `yaml:"project"`
	RepoInfo      RepoInfo      `yaml:"repo-info"`
	Policies      []Policies    `yaml:"policies"`
	Notifications Notifications `yaml:"notifications"`
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

func LoadTrustedDataToMap(filePath string, cfg *Config) {
	file, err := os.Open(filePath + TrustedDataFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// defer the closing of our file so that we can parse it later on
	defer file.Close()

	DecodeTrustedDataToMap(file, cfg)
}

func DecodeTrustedDataToMap(reader io.Reader, cfg *Config) {
	byteContent, _ := ioutil.ReadAll(reader)

	var content map[string]interface{}
	_ = yaml.Unmarshal(byteContent, &content)
	cfg.RepoInfo.TrustedData = content
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
