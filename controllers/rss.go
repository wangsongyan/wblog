package controllers

import (
	"fmt"
	"net/http"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
)

func RssGet(c *gin.Context) {
	cfg := system.GetConfiguration()
	now := helpers.GetCurrentTime()
	domain := system.GetConfiguration().Domain
	feed := &feeds.Feed{
		Title:       cfg.Title,
		Link:        &feeds.Link{Href: domain},
		Description: cfg.Seo.Description,
		Author:      &feeds.Author{Name: cfg.Seo.Author.Name, Email: cfg.Seo.Author.Email},
		Created:     now,
	}

	feed.Items = make([]*feeds.Item, 0)
	posts, err := models.ListPublishedPost("", 0, 0)
	if err != nil {
		seelog.Errorf("models.ListPublishedPost err: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, post := range posts {
		item := &feeds.Item{
			Id:          fmt.Sprintf("%s/post/%d", domain, post.ID),
			Title:       post.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/post/%d", domain, post.ID)},
			Description: string(post.Excerpt()),
			Created:     now,
		}
		feed.Items = append(feed.Items, item)
	}
	rss, err := feed.ToRss()
	if err != nil {
		seelog.Errorf("feed.ToRss err: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteString(rss)
}
