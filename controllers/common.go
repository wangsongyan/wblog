package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/cihub/seelog"

	"github.com/denisbakhtin/sitemap"
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
)

const (
	SessionKey         = "UserID"       // session key
	ContextUserKey     = "User"         // context user key
	SessionGithubState = "GITHUB_STATE" // GitHub state session key
	SessionCaptcha     = "GIN_CAPTCHA"  // captcha session key
)

func Handle404(c *gin.Context) {
	HandleMessage(c, "Sorry,I lost myself!")
}

func HandleMessage(c *gin.Context, message string) {
	c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
		"message": message,
	})
}

func sendMail(to, subject, body string) error {
	c := system.GetConfiguration()
	return helpers.SendToMail(c.SmtpUsername, c.SmtpPassword, c.SmtpHost, to, subject, body, "html")
}

func NotifyEmail(subject, body string) error {
	notifyEmailsStr := system.GetConfiguration().NotifyEmails
	if notifyEmailsStr != "" {
		notifyEmails := strings.Split(notifyEmailsStr, ";")
		emails := make([]string, 0)
		for _, email := range notifyEmails {
			if email != "" {
				emails = append(emails, email)
			}
		}
		if len(emails) > 0 {
			return sendMail(strings.Join(emails, ";"), subject, body)
		}
	}
	return nil
}

func CreateXMLSitemap() (err error) {
	configuration := system.GetConfiguration()
	folder := path.Join(configuration.Public, "sitemap")
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		seelog.Errorf("create folder:%v", err)
		return
	}
	domain := configuration.Domain
	now := helpers.GetCurrentTime()
	items := make([]sitemap.Item, 0)

	items = append(items, sitemap.Item{
		Loc:        domain,
		LastMod:    now,
		Changefreq: "daily",
		Priority:   1,
	})

	posts, err := models.ListPublishedPost("", 0, 0)
	if err != nil {
		seelog.Errorf("models.ListPublishedPost:%v", err)
		return
	}
	for _, post := range posts {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s/post/%d", domain, post.ID),
			LastMod:    post.UpdatedAt,
			Changefreq: "weekly",
			Priority:   0.9,
		})
	}

	pages, err := models.ListPublishedPage()
	if err != nil {
		seelog.Errorf("models.ListPublishedPage:%v", err)
		return
	}
	for _, page := range pages {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s/page/%d", domain, page.ID),
			LastMod:    page.UpdatedAt,
			Changefreq: "monthly",
			Priority:   0.8,
		})
	}

	err = sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), items)
	if err != nil {
		seelog.Errorf("sitemap.SiteMap:%v", err)
		return
	}
	err = sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/static/sitemap/")
	if err != nil {
		seelog.Errorf("sitemap.SiteMapIndex:%v", err)
		return
	}
	return
}

func writeJSON(ctx *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] = false
	}
	ctx.JSON(http.StatusOK, h)
}
