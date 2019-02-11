package system

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
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
	DSN                string `yaml:"dsn"`            //database dsn
	NotifyEmails       string `yaml:"notify_emails"`  //notify_emails
	PageSize           int    `yaml:"page_size"`      //page_size
	SmmsFileServer     string `yaml:"smms_fileserver"`
}

const (
	DEFAULT_PAGESIZE = 10
)

var configuration *Configuration

func LoadConfiguration(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	if config.PageSize <= 0 {
		config.PageSize = DEFAULT_PAGESIZE
	}
	configuration = &config
	return err
}

func GetConfiguration() *Configuration {
	return configuration
}
