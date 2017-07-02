package main

import (
	"github.com/gin-gonic/gin"
	//"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/controllers"
	"github.com/wangsongyan/wblog/models"
)

func main() {

	db := models.InitDB()
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLGlob("views/**/*")
	router.Static("/static", "./static")

	router.GET("/page/:id", controllers.PageGet)
	router.GET("/post/:id", controllers.PostGet)
	router.GET("/tag/:id", controllers.TagGet)

	authorized := router.Group("/admin")
	{
		// page
		authorized.GET("/page", controllers.PageIndex)
		authorized.GET("/new_page", controllers.PageNew)
		authorized.POST("/new_page", controllers.PageCreate)
		authorized.GET("/page/:id/edit", controllers.PageEdit)
		authorized.POST("/page/:id/edit", controllers.PageUpdate)
		authorized.POST("/page/:id/delete", controllers.PageDelete)

		// post
		authorized.GET("/post", controllers.PostIndex)
		authorized.GET("/new_post", controllers.PostNew)
		authorized.POST("/new_post", controllers.PostCreate)
		authorized.GET("/post/:id/edit", controllers.PostEdit)
		authorized.POST("/post/:id/edit", controllers.PostUpdate)
		authorized.POST("/post/:id/delete", controllers.PostDelete)

		// tag
		authorized.POST("/new_tag", controllers.TagCreate)
	}

	router.Run(":8090")
}
