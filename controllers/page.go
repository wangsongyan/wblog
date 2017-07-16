package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/models"
	"net/http"
	"strconv"
)

func PageGet(c *gin.Context) {
	id := c.Param("id")
	page, err := models.GetPageById(id)
	if err == nil && page.IsPublished {
		page.View++
		page.Update()
		c.HTML(http.StatusOK, "page/display.html", gin.H{
			"page": page,
		})
	} else {
		Handle404(c)
	}
}

func PageNew(c *gin.Context) {
	c.HTML(http.StatusOK, "page/new.html", nil)
}

func PageCreate(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished
	page := &models.Page{
		Title:       title,
		Body:        body,
		IsPublished: published,
	}
	err := page.Insert()
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, "/page/"+strconv.FormatUint(uint64(page.ID), 10))
	} else {
		c.HTML(http.StatusOK, "page/new.html", gin.H{
			"message": err.Error(),
			"page":    page,
		})
	}
}

func PageEdit(c *gin.Context) {
	id := c.Param("id")
	page, err := models.GetPageById(id)
	if err == nil {
		c.HTML(http.StatusOK, "page/modify.html", gin.H{
			"page": page,
		})
	} else {
		Handle404(c)
	}
}

func PageUpdate(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished
	pid, err := strconv.ParseUint(id, 10, 64)
	if err == nil {
		page := &models.Page{Title: title, Body: body, IsPublished: published}
		page.ID = uint(pid)
		err = page.Update()
		if err == nil {
			c.Redirect(http.StatusMovedPermanently, "/page/"+id)
		} else {
			// TODO
		}
	} else {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func PagePublish(c *gin.Context) {
	id := c.Param("id")
	page, err := models.GetPageById(id)
	if err == nil {
		page.IsPublished = !page.IsPublished
		err = page.Update()
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": err == nil,
	})
}

func PageDelete(c *gin.Context) {
	id := c.Param("id")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err == nil {
		page := &models.Page{}
		page.ID = uint(pid)
		page.Delete()
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"succeed": true,
			})
			return
		}
	}
	c.AbortWithError(http.StatusInternalServerError, err)
}

func PageIndex(c *gin.Context) {
	pages, _ := models.ListPage(false)
	user, _ := c.Get("User")
	c.HTML(http.StatusOK, "admin/page.html", gin.H{
		"pages": pages,
		"user":  user,
	})
}
