package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/alimoeeny/gooauth2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
	"io/ioutil"
	"log"
	"net/http"
)

func SigninGet(c *gin.Context) {
	/*session := sessions.Default(c)
	if session.Get("UserID") != nil {
		c.Redirect(http.StatusMovedPermanently, "/admin/index")
	} else {
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"authUrl": fmt.Sprintf(oauthCfg.AuthURL, oauthCfg.ClientId),
		})
	}*/
	c.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"authUrl": fmt.Sprintf(system.GetConfiguration().GithubAuthUrl, system.GetConfiguration().GithubClientId),
	})
}

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signup.html", nil)
}

func LogoutGet(c *gin.Context) {
	s := sessions.Default(c)
	s.Delete("UserID")
	s.Save()
	c.Redirect(http.StatusSeeOther, "/signin")
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
			c.HTML(http.StatusOK, "auth/signin.html", gin.H{
				"user": user,
			})
			return
		} else {
			err = errors.New("email already exists.")
		}
	}
	c.HTML(http.StatusOK, "auth/signup.html", gin.H{
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
			s := sessions.Default(c)
			s.Set("UserID", user.ID)
			s.Save()
			if user.IsAdmin {
				c.Redirect(http.StatusMovedPermanently, "/admin/index")
			} else {
				c.Redirect(http.StatusMovedPermanently, "/")
			}
			return
		} else {
			err = errors.New("invalid username or password.")
		}
	} else {
		err = errors.New("error parameter.")
	}
	c.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"message": err.Error(),
	})
}

func Oauth2Callback(c *gin.Context) {
	code := c.Query("code")
	t := &oauth.Transport{Config: &oauth.Config{
		ClientId:     system.GetConfiguration().GithubClientId,
		ClientSecret: system.GetConfiguration().GithubClientSecret,
		RedirectURL:  system.GetConfiguration().GithubRedirectURL,
		TokenURL:     system.GetConfiguration().GithubTokenUrl,
		Scope:        system.GetConfiguration().GithubScope,
	}}
	// Exchange the received code for a token
	tok, err := t.Exchange(code)
	if err == nil {
		tokenCache := oauth.CacheFile("./request.token")

		err := tokenCache.PutToken(tok)
		if err != nil {
			log.Println("Cache write:", err)
		}
		log.Printf("Token is cached in %v\n", tokenCache)
		token := tok.AccessToken
		fmt.Print(token)
		resp, err := http.Get("https://api.github.com/user?access_token=" + token)
		defer resp.Body.Close()
		if err != nil {
			log.Println(err)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			var githubinfo struct {
				Login     string `json:"login"`
				AvatarUrl string `json:"avatar_url"`
			}
			fmt.Println(string(body))
			err = json.Unmarshal(body, &githubinfo)
			if err == nil {
				fmt.Println(githubinfo)
				var user *models.User
				sessionUser, exists := c.Get("User")
				fmt.Println(sessionUser)
				if exists {
					var ok bool
					user, ok = sessionUser.(*models.User)
					if ok {
						user.GithubLoginId = githubinfo.Login
						user.AvatarUrl = githubinfo.AvatarUrl
						err = user.Update()
						//TODO 检查githubId是否已经存在
					} else {
						err = errors.New("assert failed.")
					}
				} else {
					user = &models.User{
						GithubLoginId: githubinfo.Login,
						AvatarUrl:     githubinfo.AvatarUrl,
					}
					user, err = user.FirstOrCreate()
				}

				if err == nil {
					s := sessions.Default(c)
					s.Set("UserID", user.ID)
					s.Save()
					c.Redirect(http.StatusMovedPermanently, "/admin/index")
					return
				} else {
					log.Print(err)
				}
			} else {
				log.Print(err)
			}
		}
	} else {
		log.Println(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/signin")

}

func ProfileGet(c *gin.Context) {
	sessionUser, exists := c.Get("User")
	if exists {
		c.HTML(http.StatusOK, "admin/profile.html", gin.H{
			"user":    sessionUser,
			"authUrl": fmt.Sprintf(system.GetConfiguration().GithubAuthUrl, system.GetConfiguration().GithubClientId),
		})
	}
}

func ProfileUpdate(c *gin.Context) {
	avatarUrl := c.PostForm("avatarUrl")
	nickName := c.PostForm("nickName")
	sessionUser, _ := c.Get("User")
	if user, ok := sessionUser.(*models.User); ok {
		err := user.UpdateProfile(avatarUrl, nickName)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"succeed": true,
				"user":    models.User{AvatarUrl: avatarUrl, NickName: nickName},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
				"message": err.Error(),
			})
		}
	}
}

func BindEmail(c *gin.Context) {
	email := c.PostForm("email")
	sessionUser, _ := c.Get("User")
	if user, ok := sessionUser.(*models.User); ok {
		if len(user.Email) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
				"message": "email have bound.",
			})
		} else {
			_, err := models.GetUserByUsername(email)
			if err != nil {
				err := user.UpdateEmail(email)
				c.JSON(http.StatusOK, gin.H{
					"succeed": err == nil,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"succeed": false,
					"message": "email have be registered!",
				})
			}
		}
	}
}

func UnbindEmail(c *gin.Context) {
	sessionUser, _ := c.Get("User")
	if user, ok := sessionUser.(*models.User); ok {
		if len(user.Email) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
				"message": "email haven't bound.",
			})
		} else {
			err := user.UpdateEmail("")
			c.JSON(http.StatusOK, gin.H{
				"succeed": err == nil,
			})
		}
	}
}

func UnbindGithub(c *gin.Context) {
	sessionUser, _ := c.Get("User")
	if user, ok := sessionUser.(*models.User); ok {
		if len(user.GithubLoginId) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
				"message": "github haven't bound.",
			})
		} else {
			err := user.UpdateGithubId("")
			c.JSON(http.StatusOK, gin.H{
				"succeed": err == nil,
			})
		}
	}
}

func UserIndex(c *gin.Context) {
	users, _ := models.ListUsers()
	user, _ := c.Get("User")
	c.HTML(http.StatusOK, "admin/user.html", gin.H{
		"users": users,
		"user":  user,
	})
}
