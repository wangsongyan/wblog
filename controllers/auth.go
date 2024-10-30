package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/system"
)

func AuthGet(c *gin.Context) {
	authType := c.Param("authType")

	session := sessions.Default(c)
	uuid := helpers.UUID()
	session.Delete(SessionGithubState)
	session.Set(SessionGithubState, uuid)
	session.Save()

	cfg := system.GetConfiguration()

	authurl := "/signin"
	switch authType {
	case "github":
		if cfg.Github.Enabled {
			authurl = fmt.Sprintf(cfg.Github.AuthUrl, cfg.Github.ClientId, uuid)
		}
	case "weibo":
	case "qq":
	case "wechat":
	case "oschina":
	default:
	}
	c.Redirect(http.StatusFound, authurl)
}
