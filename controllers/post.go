package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/models"
	"net/http"
	"strconv"
)

func PostGet(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err == nil {
		post.Tags, _ = models.ListTagByPostId(id)
		post.Comments, _ = models.ListCommentByPostID(id)
		c.HTML(http.StatusOK, "post/display.html", gin.H{
			"post": post,
		})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PostNew(c *gin.Context) {
	c.HTML(http.StatusOK, "post/new.html", nil)
}

func PostCreate(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")

	post := &models.Post{
		Title:       title,
		Body:        body,
		IsPublished: isPublished,
	}
	err := post.Insert()
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, "/post/"+strconv.FormatUint(uint64(post.ID), 10))
	} else {
		c.HTML(http.StatusOK, "post/new.html", gin.H{
			"post":    post,
			"message": err.Error(),
		})
	}
}

func PostEdit(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err == nil {
		tags, _ := models.ListTagByPostId(id)
		c.HTML(http.StatusOK, "post/modify.html", gin.H{
			"post": post,
			"tags": tags,
		})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PostUpdate(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err == nil {
		post := &models.Post{
			Title:       title,
			Body:        body,
			IsPublished: isPublished,
		}
		post.ID = uint(pid)
		err = post.Update()
		if err == nil {
			c.Redirect(http.StatusMovedPermanently, "/post/"+id)
		} else {
			c.HTML(http.StatusOK, "post/modify.html", gin.H{
				"post":    post,
				"message": err.Error(),
			})
		}
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PostDelete(c *gin.Context) {
	id := c.PostForm("id")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err == nil {
		post := &models.Post{}
		post.ID = uint(pid)
		post.Delete()
		c.Redirect(http.StatusMovedPermanently, "/admin/post")
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PostIndex(c *gin.Context) {
	posts, err := models.ListPost("")
	c.HTML(http.StatusOK, "", gin.H{
		"posts":   posts,
		"message": err.Error(),
	})
}
