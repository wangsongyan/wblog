package system

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
	"os"
)

type (
	Backup struct {
		Enabled   bool   `toml:"enabled"`
		BackupKey string `toml:"backup_key"`
	}

	Database struct {
		Dialect string `toml:"dialect"`
		DSN     string `toml:"dsn"`
	}

	Author struct {
		Name  string `toml:"name"`
		Email string `toml:"email"`
	}

	Seo struct {
		Description string `toml:"description"`
		Author      Author `toml:"author"`
	}

	Qiniu struct {
		Enabled    bool   `toml:"enabled"`
		AccessKey  string `toml:"accesskey"`
		SecretKey  string `toml:"secretkey"`
		FileServer string `toml:"fileserver"`
		Bucket     string `toml:"bucket"`
	}

	Smms struct {
		Enabled bool   `toml:"enabled"`
		ApiUrl  string `toml:"apiurl"`
		ApiKey  string `toml:"apikey"`
	}

	Github struct {
		Enabled      bool   `toml:"enabled"`
		ClientId     string `toml:"clientid"`
		ClientSecret string `toml:"clientsecret"`
		RedirectURL  string `toml:"redirecturl"`
		AuthUrl      string `toml:"authurl"`
		TokenUrl     string `toml:"tokenurl"`
		Scope        string `toml:"scope"`
	}

	Smtp struct {
		Enabled  bool   `toml:"enabled"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Host     string `toml:"host"`
	}

	Navigator struct {
		Title  string `toml:"title"`
		Url    string `toml:"url"`
		Target string `toml:"target"`
	}

	Configuration struct {
		Addr          string      `toml:"addr"`
		SignupEnabled bool        `toml:"signup_enabled"`
		Title         string      `toml:"title"`
		SessionSecret string      `toml:"session_secret"`
		Domain        string      `toml:"domain"`
		FileServer    string      `toml:"file_server"`
		NotifyEmails  string      `toml:"notify_emails"`
		PageSize      int         `toml:"page_size"`
		PublicDir     string      `toml:"public"`
		ViewDir       string      `toml:"view"`
		Database      Database    `toml:"database"`
		Seo           Seo         `toml:"seo"`
		Qiniu         Qiniu       `toml:"qiniu"`
		Smms          Smms        `toml:"smms"`
		Github        Github      `toml:"github"`
		Smtp          Smtp        `toml:"smtp"`
		Navigators    []Navigator `toml:"navigators"`
		Backup        Backup      `toml:"backup"`
	}
)

func (a Author) String() string {
	return fmt.Sprintf("%s,%s", a.Name, a.Email)
}

var configuration *Configuration

func defaultConfig() Configuration {
	return Configuration{
		Addr:          ":8090",
		Domain:        "https://example.com",
		Title:         "Wblog",
		SessionSecret: "wblog",
		FileServer:    "smms",
		PageSize:      10,
		PublicDir:     "static",
		ViewDir:       "views/**/*",
		Database: Database{
			Dialect: "sqlite",
			DSN:     "wblog.db?_loc=Asia/Shanghai",
		},
		Seo: Seo{
			Description: "Wblog,talk about golang,java and so on.",
			Author: Author{
				Name:  "wangsy",
				Email: "wangsy0129@qq.com",
			},
		},
		Qiniu: Qiniu{
			AccessKey:  "",
			SecretKey:  "",
			FileServer: "",
			Bucket:     "wblog",
		},
		Smms: Smms{
			ApiUrl: "https://sm.ms/api/v2/upload",
		},
		Github: Github{
			ClientId:     "",
			ClientSecret: "",
			RedirectURL:  "https://example.com/oauth2callback",
			AuthUrl:      "https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&state=%s",
			TokenUrl:     "https://github.com/login/oauth/access_token",
			Scope:        "",
		},
		Smtp: Smtp{
			Username: "",
			Password: "",
			Host:     "smtp.163.com:25",
		},
		Navigators: []Navigator{
			{
				Title: "Posts",
				Url:   "/index",
			},
			{
				Title: "AboutMe",
				Url:   "/page/6",
			},
			{
				Title:  "RSS",
				Url:    "/rss",
				Target: "_blank",
			},
			{
				Title: "Subscribe",
				Url:   "/subscribe",
			},
		},
	}
}

func LoadConfiguration(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var config = defaultConfig()
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	configuration = &config
	return nil
}

func Generate() error {
	config := defaultConfig()
	placeholder := "[!!]"
	config.Domain = placeholder
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile("conf/conf.sample.toml", data, os.ModePerm)
}

func GetConfiguration() *Configuration {
	return configuration
}
