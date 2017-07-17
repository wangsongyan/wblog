package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/claudiu/gocron"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/controllers"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
	"html/template"
	"net/http"
)

func main() {

	db := models.InitDB()
	defer db.Close()
	system.LoadConfiguration("conf/conf.yaml")

	router := gin.Default()

	//router.LoadHTMLGlob("views/**/*")

	setTemplate(router)
	setSessions(router)
	router.Use(SharedData())

	//Periodic tasks
	gocron.Every(1).Day().Do(system.CreateXMLSitemap)
	gocron.Start()

	router.Static("/static", "./static")

	router.NoRoute(controllers.Handle404)
	router.GET("/", controllers.IndexGet)
	router.GET("/index", controllers.IndexGet)
	router.GET("/rss", controllers.RssGet)

	if system.GetConfiguration().SignupEnabled {
		router.GET("/signup", controllers.SignupGet)
		router.POST("/signup", controllers.SignupPost)
	}
	// user signin and logout
	router.GET("/signin", controllers.SigninGet)
	router.POST("/signin", controllers.SigninPost)
	router.GET("/logout", controllers.LogoutGet)
	router.GET("/oauth2callback", controllers.Oauth2Callback)

	router.POST("/new_comment", controllers.CommentPost)
	router.POST("/comment/:id/delete", controllers.CommentDelete)

	// subscriber
	router.GET("/subscribe", controllers.SubscribeGet)
	router.POST("/subscribe", controllers.Subscribe)
	router.GET("/active", controllers.ActiveSubsciber)
	router.GET("/unsubscribe", controllers.UnSubscribe)

	router.GET("/page/:id", controllers.PageGet)
	router.GET("/post/:id", controllers.PostGet)
	router.GET("/tag/:id", controllers.TagGet)
	router.GET("/archives/:year/:month", controllers.ArchiveGet)

	router.GET("/link/:id", controllers.LinkGet)

	authorized := router.Group("/admin")
	authorized.Use(AuthRequired())
	{
		// index
		authorized.GET("/index", controllers.AdminIndex)

		// image upload
		authorized.POST("/upload", controllers.Upload)

		// page
		authorized.GET("/page", controllers.PageIndex)
		authorized.GET("/new_page", controllers.PageNew)
		authorized.POST("/new_page", controllers.PageCreate)
		authorized.GET("/page/:id/edit", controllers.PageEdit)
		authorized.POST("/page/:id/edit", controllers.PageUpdate)
		authorized.POST("/page/:id/publish", controllers.PagePublish)
		authorized.POST("/page/:id/delete", controllers.PageDelete)

		// post
		authorized.GET("/post", controllers.PostIndex)
		authorized.GET("/new_post", controllers.PostNew)
		authorized.POST("/new_post", controllers.PostCreate)
		authorized.GET("/post/:id/edit", controllers.PostEdit)
		authorized.POST("/post/:id/edit", controllers.PostUpdate)
		authorized.POST("/post/:id/publish", controllers.PostPublish)
		authorized.POST("/post/:id/delete", controllers.PostDelete)

		// tag
		authorized.POST("/new_tag", controllers.TagCreate)

		//
		authorized.GET("/user", controllers.UserIndex)

		// profile
		authorized.GET("/profile", controllers.ProfileGet)
		authorized.POST("/profile", controllers.ProfileUpdate)
		authorized.POST("/profile/email/bind", controllers.BindEmail)
		authorized.POST("/profile/email/unbind", controllers.UnbindEmail)
		authorized.POST("/profile/github/unbind", controllers.UnbindGithub)

		// subscriber
		authorized.GET("/subscriber", controllers.SubscriberIndex)
		authorized.POST("/subscriber", controllers.SubscriberPost)

		// link
		authorized.GET("/link", controllers.LinkIndex)
		authorized.POST("/new_link", controllers.LinkCreate)
		authorized.POST("/link/:id/edit", controllers.LinkUpdate)
		authorized.POST("/link/:id/delete", controllers.LinkDelete)
	}

	router.Run(system.GetConfiguration().Addr)
}

func setTemplate(engine *gin.Engine) {

	funcMap := template.FuncMap{
		"dateFormat": helpers.DateFormat,
		"substring":  helpers.Substring,
		"isOdd":      helpers.IsOdd,
		"isEven":     helpers.IsEven,
		"truncate":   helpers.Truncate,
		"add":        helpers.Add,
	}

	if gin.IsDebugging() {
		render := helpers.New()
		render.FuncMap = funcMap
		render.Glob = "views/**/*"
		engine.HTMLRender = render
	} else {
		t, err := template.ParseGlob("views/**/*")
		if err == nil {
			t.Funcs(funcMap)
		}
		engine.SetHTMLTemplate(template.Must(t, err))
	}

}

//setSessions initializes sessions & csrf middlewares
func setSessions(router *gin.Engine) {
	config := system.GetConfiguration()
	//https://github.com/gin-gonic/contrib/tree/master/sessions
	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
	//https://github.com/utrack/gin-csrf
	/*router.Use(csrf.Middleware(csrf.Options{
		Secret: config.SessionSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))*/
}

//+++++++++++++ middlewares +++++++++++++++++++++++

//SharedData fills in common data, such as user info, etc...
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if uID := session.Get("UserID"); uID != nil {
			user, _ := models.GetUser(uID)
			if user.ID != 0 {
				c.Set("User", user)
			}
		}
		if system.GetConfiguration().SignupEnabled {
			c.Set("SignupEnabled", true)
		}
		c.Next()
	}
}

//AuthRequired grants access to authenticated users, requires SharedData middleware
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("User"); user != nil {
			if u, ok := user.(*models.User); ok && u.IsAdmin {
				c.Next()
				return
			}
		}
		logrus.Warnf("User not authorized to visit %s", c.Request.RequestURI)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "Forbidden!",
		})
		c.Abort()
	}
}
