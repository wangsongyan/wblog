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
	if err == nil {
		c.HTML(http.StatusOK, "page/display.html", gin.H{
			"page": page,
		})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PageNew(c *gin.Context) {
	c.HTML(http.StatusOK, "page/new.html", nil)
}

func PageCreate(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	page := &models.Page{
		Title:       title,
		Body:        body,
		IsPublished: isPublished,
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
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PageUpdate(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err == nil {
		page := &models.Page{Title: title, Body: body, IsPublished: isPublished}
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

func PageDelete(c *gin.Context) {
	id := c.PostForm("id")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err == nil {
		page := &models.Page{}
		page.ID = uint(pid)
		page.Delete()
		if err == nil {
			//c.Redirect(http.StatusMovedPermanently, "/page/"+id)
		} else {
			// TODO
		}
	} else {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func PageIndex(c *gin.Context) {
	pages, _ := models.ListPage()
	c.HTML(http.StatusOK, "page/", gin.H{
		"pages": pages,
	})
}
