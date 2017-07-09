package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handle404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
		"message": "Sorry,I lost myself!",
	})
}
