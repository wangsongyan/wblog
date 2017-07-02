package system

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

type Configuration struct {
	SignupEnabled   bool   `yaml:"signup_enabled"` // signup enabled or not
	QiniuAppKey     string `yaml:"qiniu_appkey"`   // qiniu
	QiniuAppSecret  string `yaml:"qiniu_appsecret"`
	GithubAppKey    string `yaml:"github_appkey"` // github
	GithubAppSecret string `yaml:"github_appsecret"`
	SmtpUsername    string `yaml:"smtp_username"`  // username
	SmtpPassword    string `yaml:"smtp_password"`  //password
	SmtpHost        string `yaml:"smtp_host"`      //host
	SessionSecret   string `yaml:"session_secret"` //session_secret
}

var configuration *Configuration

func LoadConfiguration(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	configuration = &config
	log.Println(configuration)
}

func GetConfiguration() *Configuration {
	return configuration
}
