package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/models"
	"net/http"
)

func SigninGet(c *gin.Context) {
	c.HTML(http.StatusOK, "user/signin.html", nil)
}

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "user/signup.html", nil)
}

func LogoutGet(c *gin.Context) {

}

func SignupPost(c *gin.Context) {
	email := c.PostForm("email")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	user := &models.User{
		Email:     email,
		Telephone: telephone,
		Password:  password,
	}
	var err error
	if len(user.Email) == 0 /*|| len(user.Telephone) == 0 */ || len(user.Password) == 0 {
		err = errors.New("error parameter.")
	} else {
		user.Password = helpers.Md5(user.Email + user.Password)
		err = user.Insert()
		if err == nil {
			c.HTML(http.StatusOK, "user/signin.html", gin.H{
				"user": user,
			})
			return
		} else {
			err = errors.New("email already exists.")
		}
	}
	c.HTML(http.StatusOK, "user/signup.html", gin.H{
		"message": err.Error(),
	})
}

func SigninPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	var err error
	if len(username) > 0 && len(password) > 0 {
		var user *models.User
		user, err = models.GetUserByUsername(username)
		if err == nil && user.Password == helpers.Md5(username+password) {
			c.Redirect(http.StatusMovedPermanently, "/admin/index")
			return
		} else {
			err = errors.New("invalid username or password.")
		}
	} else {
		err = errors.New("error parameter.")
	}
	c.HTML(http.StatusOK, "user/signin.html", gin.H{
		"message": err.Error(),
	})
}
