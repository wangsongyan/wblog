package system

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type Configuration struct {
	SignupEnabled      bool   `yaml:"signup_enabled"`  // signup enabled or not
	QiniuAccessKey     string `yaml:"qiniu_accesskey"` // qiniu
	QiniuSecretKey     string `yaml:"qiniu_secretkey"`
	QiniuFileServer    string `yaml:"qiniu_fileserver"`
	QiniuBucket        string `yaml:"qiniu_bucket"`
	GithubClientId     string `yaml:"github_clientid"` // github
	GithubClientSecret string `yaml:"github_clientsecret"`
	GithubAuthUrl      string `yaml:"github_authurl"`
	GithubRedirectURL  string `yaml:"github_redirecturl"`
	GithubTokenUrl     string `yaml:"github_tokenurl"`
	GithubScope        string `yaml:"github_scope"`
	SmtpUsername       string `yaml:"smtp_username"`  // username
	SmtpPassword       string `yaml:"smtp_password"`  //password
	SmtpHost           string `yaml:"smtp_host"`      //host
	SessionSecret      string `yaml:"session_secret"` //session_secret
	Domain             string `yaml:"domain"`         //domain
	Public             string `yaml:"public"`         //public
	Addr               string `yaml:"addr"`           //addr
	BackupKey          string `yaml:"backup_key"`     //backup_key
}

var configuration *Configuration

func LoadConfiguration(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	configuration = &config
	logrus.Debugf("config load succeessfully:%v", configuration)
}

func GetConfiguration() *Configuration {
	return configuration
}
