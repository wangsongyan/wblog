package controllers

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
	"net/http"
)

func PageGet(c *gin.Context) {
	id, err := ParamUint(c, "id")
	if err != nil {
		HandleMessage(c, err.Error())
		return
	}
	page, err := models.GetPageById(id)
	if err != nil || !page.IsPublished {
		Handle404(c)
		return
	}
	page.View++
	page.UpdateView()
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "page/display.html", gin.H{
		"page": page,
		"user": user,
		"cfg":  system.GetConfiguration(),
	})
}

func PageNew(c *gin.Context) {
	c.HTML(http.StatusOK, "page/new.html", gin.H{
		"user": c.MustGet(ContextUserKey),
		"cfg":  system.GetConfiguration(),
	})
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
	if err != nil {
		c.HTML(http.StatusOK, "page/new.html", gin.H{
			"message": err.Error(),
			"page":    page,
			"user":    c.MustGet(ContextUserKey),
			"cfg":     system.GetConfiguration(),
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/page")
}

func PageEdit(c *gin.Context) {
	id, err := ParamUint(c, "id")
	if err != nil {
		HandleMessage(c, err.Error())
		return
	}
	page, err := models.GetPageById(id)
	if err != nil {
		Handle404(c)
		return
	}
	c.HTML(http.StatusOK, "page/modify.html", gin.H{
		"page": page,
		"user": c.MustGet(ContextUserKey),
		"cfg":  system.GetConfiguration(),
	})
}

func PageUpdate(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished
	id, err := ParamUint(c, "id")
	if err != nil {
		HandleMessage(c, err.Error())
		return
	}
	page := &models.Page{Title: title, Body: body, IsPublished: published}
	page.ID = id
	err = page.Update()
	if err != nil {
		seelog.Errorf("page.Update err: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/page")
}

func PagePublish(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id, err := ParamUint(c, "id")
	if err != nil {
		HandleMessage(c, err.Error())
		return
	}
	page, err := models.GetPageById(id)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	page.IsPublished = !page.IsPublished
	err = page.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func PageDelete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id, err := ParamUint(c, "id")
	if err != nil {
		res["message"] = err.Error()
		return
	}
	page := &models.Page{}
	page.ID = id
	err = page.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func PageIndex(c *gin.Context) {
	pages, _ := models.ListAllPage()
	c.HTML(http.StatusOK, "admin/page.html", gin.H{
		"pages":    pages,
		"user":     c.MustGet(ContextUserKey),
		"comments": models.MustListUnreadComment(),
		"cfg":      system.GetConfiguration(),
	})
}
