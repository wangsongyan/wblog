package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/models"
	"net/http"
)

func TagCreate(c *gin.Context) {
	name := c.PostForm("name")
	tag := &models.Tag{Name: name}
	err := tag.Insert()
	if err != nil {
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
		c.HTML(http.StatusOK, "", gin.H{
			"posts": posts,
		})
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
