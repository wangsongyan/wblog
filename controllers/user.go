package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/alimoeeny/gooauth2"
	"github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
	"io/ioutil"
	"net/http"
	"strconv"
)

type GithubUserInfo struct {
	AvatarURL         string      `json:"avatar_url"`
	Bio               interface{} `json:"bio"`
	Blog              string      `json:"blog"`
	Company           interface{} `json:"company"`
	CreatedAt         string      `json:"created_at"`
	Email             interface{} `json:"email"`
	EventsURL         string      `json:"events_url"`
	Followers         int         `json:"followers"`
	FollowersURL      string      `json:"followers_url"`
	Following         int         `json:"following"`
	FollowingURL      string      `json:"following_url"`
	GistsURL          string      `json:"gists_url"`
	GravatarID        string      `json:"gravatar_id"`
	Hireable          interface{} `json:"hireable"`
	HTMLURL           string      `json:"html_url"`
	ID                int         `json:"id"`
	Location          interface{} `json:"location"`
	Login             string      `json:"login"`
	Name              interface{} `json:"name"`
	OrganizationsURL  string      `json:"organizations_url"`
	PublicGists       int         `json:"public_gists"`
	PublicRepos       int         `json:"public_repos"`
	ReceivedEventsURL string      `json:"received_events_url"`
	ReposURL          string      `json:"repos_url"`
	SiteAdmin         bool        `json:"site_admin"`
	StarredURL        string      `json:"starred_url"`
	SubscriptionsURL  string      `json:"subscriptions_url"`
	Type              string      `json:"type"`
	UpdatedAt         string      `json:"updated_at"`
	URL               string      `json:"url"`
}

func SigninGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signin.html", nil)
}

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signup.html", nil)
}

func LogoutGet(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
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
		IsAdmin:   true,
	}
	var err error
	if len(user.Email) == 0 || len(user.Password) == 0 {
		err = errors.New("email or password cannot be null.")
	} else {
		user.Password = helpers.Md5(user.Email + user.Password)
		err = user.Insert()
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"succeed": true,
			})
			return
		} else {
			err = errors.New("email already exists.")
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": false,
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
		fmt.Println(user, err)
		if err == nil && user.Password == helpers.Md5(username+password) {
			if !user.LockState {
				s := sessions.Default(c)
				s.Clear()
				s.Set(SESSION_KEY, user.ID)
				s.Save()
				if user.IsAdmin {
					c.Redirect(http.StatusMovedPermanently, "/admin/index")
				} else {
					c.Redirect(http.StatusMovedPermanently, "/")
				}
				return
			} else {
				err = errors.New("Your account have been locked.")
			}
		} else {
			err = errors.New("invalid username or password.")
		}
	} else {
		err = errors.New("username or password cannot be null.")
	}
	c.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"message": err.Error(),
	})
}

func Oauth2Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	session := sessions.Default(c)
	if len(state) == 0 || state != session.Get(SESSION_GITHUB_STATE) {
		c.Abort()
		return
	} else {
		session.Delete(SESSION_GITHUB_STATE)
		session.Save()
	}
	token, err := exchangeTokenByCode(code)
	if err == nil {
		var userInfo *GithubUserInfo
		userInfo, err = getGithubUserInfoByAceessToken(token)
		if err == nil {
			var user *models.User
			if sessionUser, exists := c.Get(CONTEXT_USER_KEY); exists {
				user, _ = sessionUser.(*models.User)
				_, err1 := models.IsGithubIdExists(userInfo.Login, user.ID)
				if err1 != nil { // 未绑定
					if user.IsAdmin {
						user.GithubLoginId = userInfo.Login
					}
					user.AvatarUrl = userInfo.AvatarURL
					user.GithubUrl = userInfo.HTMLURL
					err = user.UpdateGithubUserInfo()
				} else {
					err = errors.New("this github loginId has bound another account.")
				}
			} else {
				user = &models.User{
					GithubLoginId: userInfo.Login,
					AvatarUrl:     userInfo.AvatarURL,
					GithubUrl:     userInfo.HTMLURL,
				}
				user, err = user.FirstOrCreate()
				if err == nil {
					if user.LockState {
						err = errors.New("Your account have been locked.")
						HandleMessage(c, "Your account have been locked.")
						return
					}
				}
			}

			if err == nil {
				s := sessions.Default(c)
				s.Clear()
				s.Set(SESSION_KEY, user.ID)
				s.Save()
				if user.IsAdmin {
					c.Redirect(http.StatusMovedPermanently, "/admin/index")
				} else {
					c.Redirect(http.StatusMovedPermanently, "/")
				}
				return
			}
		}
	}
	seelog.Error(err)
	c.Redirect(http.StatusMovedPermanently, "/signin")
}

func exchangeTokenByCode(code string) (string, error) {
	t := &oauth.Transport{Config: &oauth.Config{
		ClientId:     system.GetConfiguration().GithubClientId,
		ClientSecret: system.GetConfiguration().GithubClientSecret,
		RedirectURL:  system.GetConfiguration().GithubRedirectURL,
		TokenURL:     system.GetConfiguration().GithubTokenUrl,
		Scope:        system.GetConfiguration().GithubScope,
	}}
	tok, err := t.Exchange(code)
	if err == nil {
		tokenCache := oauth.CacheFile("./request.token")
		err := tokenCache.PutToken(tok)
		return tok.AccessToken, err
	}
	return "", err
}

func getGithubUserInfoByAceessToken(token string) (*GithubUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/user?access_token=%s", token))
	defer resp.Body.Close()
	if err == nil {
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			var userInfo GithubUserInfo
			err = json.Unmarshal(body, &userInfo)
			return &userInfo, err
		}
	}
	return nil, err
}

func ProfileGet(c *gin.Context) {
	sessionUser, exists := c.Get(CONTEXT_USER_KEY)
	if exists {
		c.HTML(http.StatusOK, "admin/profile.html", gin.H{
			"user":     sessionUser,
			"comments": models.MustListUnreadComment(),
		})
	}
}

func ProfileUpdate(c *gin.Context) {
	avatarUrl := c.PostForm("avatarUrl")
	nickName := c.PostForm("nickName")
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
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
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
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
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
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
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
	if user, ok := sessionUser.(*models.User); ok {
		if len(user.GithubLoginId) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
				"message": "github haven't bound.",
			})
		} else {
			user.GithubLoginId = ""
			err := user.UpdateGithubUserInfo()
			c.JSON(http.StatusOK, gin.H{
				"succeed": err == nil,
			})
		}
	}
}

func UserIndex(c *gin.Context) {
	users, _ := models.ListUsers()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/user.html", gin.H{
		"users":    users,
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}

func UserLock(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.ParseUint(id, 10, 64)
	user, err := models.GetUser(uint(_id))
	if err == nil {
		user.LockState = !user.LockState
		err = user.Lock()
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
