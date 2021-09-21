package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Project struct {
	Owner string 					`yaml:"owner"`
	Repo string 					`yaml:"repo"`
}

type RepoInfoChecks struct {
	TrustedDataFile string 			`yaml:"trusted-data-file"`
	CiCdPath string 				`yaml:"ci-cd-path"`
	ProtectedBranches []string 		`yaml:"protected-branches"`
}

type Config struct {
	Project  Project        		`yaml:"project"`
	RepoInfoChecks RepoInfoChecks 	`yaml:"repo-info-checks"`
}

func LoadConfig(filename string, cfg *Config)  {
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

func LoadFileToJsonMap(filename string) map[string]interface{} {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteContent, _ := ioutil.ReadAll(jsonFile)

	var content map[string]interface{}
	_ = json.Unmarshal(byteContent, &content)

	return content
}