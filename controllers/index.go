package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/wangsongyan/wblog/models"
	"net/http"
	"strconv"
)

func IndexGet(c *gin.Context) {
	posts, err := models.ListPublishedPost("")
	if err == nil {
		policy := bluemonday.StrictPolicy()
		for _, post := range posts {
			post.Tags, _ = models.ListTagByPostId(strconv.FormatUint(uint64(post.ID), 10))
			post.Body = policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Body))))
		}
		c.HTML(http.StatusOK, "index/index.html", gin.H{
			"posts":    posts,
			"tags":     models.MustListTag(),
			"archives": models.MustListPostArchives(),
		})
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AdminIndex(c *gin.Context) {
	user, _ := c.Get("User")
	c.HTML(http.StatusOK, "admin/index.html", gin.H{
		"pageCount":    models.CountPage(),
		"postCount":    models.CountPost(),
		"tagCount":     models.CountTag(),
		"commentCount": models.CountComment(),
		"user":         user,
	})
}
