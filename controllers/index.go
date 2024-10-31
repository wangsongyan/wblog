package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
)

func IndexGet(c *gin.Context) {
	var (
		pageIndex int
		pageSize  = system.GetConfiguration().PageSize
		total     int
		page      string
		err       error
		posts     []*models.Post
		policy    *bluemonday.Policy
	)
	page = c.Query("page")
	pageIndex, _ = strconv.Atoi(page)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	posts, err = models.ListPublishedPost("", pageIndex, pageSize)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	total, err = models.CountPostByTag("")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Tags, _ = models.ListTagByPostId(strconv.FormatUint(uint64(post.ID), 10))
		post.Body = policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Body))))
	}
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"posts":           posts,
		"tags":            models.MustListTag(),
		"archives":        models.MustListPostArchives(),
		"links":           models.MustListLinks(),
		"user":            user,
		"pageIndex":       pageIndex,
		"totalPage":       int(math.Ceil(float64(total) / float64(pageSize))),
		"path":            c.Request.URL.Path,
		"maxReadPosts":    models.MustListMaxReadPost(),
		"maxCommentPosts": models.MustListMaxCommentPost(),
		"cfg":             system.GetConfiguration(),
	})
}

func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/index.html", gin.H{
		"pageCount":    models.CountPage(),
		"postCount":    models.CountPost(),
		"tagCount":     models.CountTag(),
		"commentCount": models.CountComment(),
		"user":         c.MustGet(ContextUserKey),
		"comments":     models.MustListUnreadComment(),
		"cfg":          system.GetConfiguration(),
	})
}
