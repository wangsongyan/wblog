package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	var (
		err      error
		url      string
		uploader Uploader
	)
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"succeed": false,
			"message": err.Error(),
		})
		return
	}

	//uploader = QiniuUploader{}
	uploader = SmmsUploader{}

	url, err = uploader.upload(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"succeed": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": true,
		"url":     url,
	})
}
