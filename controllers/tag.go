package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/wangsongyan/wblog/models"
	"net/http"
	"strconv"
)

func TagCreate(c *gin.Context) {
	name := c.PostForm("value")
	tag := &models.Tag{Name: name}
	err := tag.Insert()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"data": tag,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	}
}

func TagGet(c *gin.Context) {
	id := c.Param("id")
	posts, err := models.ListPost(id)
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
