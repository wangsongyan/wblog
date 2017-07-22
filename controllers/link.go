package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wangsongyan/wblog/models"
	"net/http"
	"strconv"
)

func LinkIndex(c *gin.Context) {
	links, _ := models.ListLinks()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/link.html", gin.H{
		"links":    links,
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}

func LinkCreate(c *gin.Context) {
	name := c.PostForm("name")
	url := c.PostForm("url")
	sort := c.PostForm("sort")
	if len(name) == 0 || len(url) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"succeed": false,
			"message": "error parameter!",
		})
	} else {
		_sort, err := strconv.ParseInt(sort, 10, 64)
		if err == nil {
			link := &models.Link{
				Name: name,
				Url:  url,
				Sort: int(_sort),
			}
			err = link.Insert()
		}
		c.JSON(http.StatusOK, gin.H{
			"succeed": err == nil,
		})
	}
}

func LinkUpdate(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	url := c.PostForm("url")
	sort := c.PostForm("sort")
	if len(id) == 0 || len(name) == 0 || len(url) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"succeed": false,
			"message": "error parameter!",
		})
	} else {
		var err error
		_id, err := strconv.ParseUint(id, 10, 64)
		_sort, err := strconv.ParseInt(sort, 10, 64)
		if err == nil {
			link := &models.Link{
				Name: name,
				Url:  url,
				Sort: int(_sort),
			}
			link.ID = uint(_id)
			err = link.Update()
		}

		c.JSON(http.StatusOK, gin.H{
			"succeed": err == nil,
		})
	}
}

func LinkGet(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.ParseInt(id, 10, 64)
	link, err := models.GetLinkById(uint(_id))
	if err == nil {
		link.View++
		link.Update()
		c.Redirect(http.StatusFound, link.Url)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func LinkDelete(c *gin.Context) {
	id := c.Param("id")
	var err error
	_id, err := strconv.ParseUint(id, 10, 64)
	if err == nil {
		link := new(models.Link)
		link.ID = uint(_id)
		err = link.Delete()
	}
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"succeed": true,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"succeed": false,
			"message": err.Error(),
		})
	}
}
