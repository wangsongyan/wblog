package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/models"
	"net/http"
	"strconv"
)

func LinkIndex(c *gin.Context) {
	links, _ := models.ListLinks()
	user, _ := c.Get("User")
	c.HTML(http.StatusOK, "admin/link.html", gin.H{
		"links": links,
		"user":  user,
	})
}

func LinkCreate(c *gin.Context) {
	name := c.PostForm("name")
	url := c.PostForm("url")
	if len(name) == 0 || len(url) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"succeed": false,
			"message": "error parameter!",
		})
	} else {
		link := &models.Link{
			Name: name,
			Url:  url,
		}
		err := link.Insert()
		c.JSON(http.StatusOK, gin.H{
			"succeed": err == nil,
		})
	}
}

func LinkUpdate(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	url := c.PostForm("url")
	if len(id) == 0 || len(name) == 0 || len(url) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"succeed": false,
			"message": "error parameter!",
		})
	} else {
		var err error
		_id, err := strconv.ParseUint(id, 10, 64)
		if err == nil {
			link := &models.Link{
				Name: name,
				Url:  url,
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
	url := c.Query("url")
	c.Redirect(http.StatusMovedPermanently, url)
}
