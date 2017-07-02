package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/models"
	"net/http"
	"strconv"
)

func IndexGet(c *gin.Context) {
	posts, err := models.ListPost("")
	if err == nil {
		for _, post := range posts {
			post.Tags, _ = models.ListTagByPostId(strconv.FormatUint(uint64(post.ID), 10))
		}
		c.HTML(http.StatusOK, "index/index.html", gin.H{
			"posts":    posts,
			"tags":     helpers.ListTag(),
			"archives": helpers.ListArchive(),
		})
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
