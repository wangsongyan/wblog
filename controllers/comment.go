package controllers

import (
	"net/http"
	"strconv"

	"fmt"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
)

func CommentPost(c *gin.Context) {
	s := sessions.Default(c)
	sessionUserID := s.Get(SESSION_KEY)
	userId, _ := sessionUserID.(uint)

	verifyCode := c.PostForm("verifyCode")
	captchaId := s.Get(SESSION_CAPTCHA)
	s.Delete(SESSION_CAPTCHA)
	_captchaId, _ := captchaId.(string)
	if !captcha.VerifyString(_captchaId, verifyCode) {
		c.JSON(http.StatusOK, gin.H{
			"succeed": false,
			"message": "error verifycode",
		})
		return
	}

	var err error
	postId := c.PostForm("postId")
	content := c.PostForm("content")
	if len(content) == 0 {
		err = errors.New("content cannot be empty.")
	}
	var post *models.Post
	post, err = models.GetPostById(postId)
	if err == nil {
		pid, err := strconv.ParseUint(postId, 10, 64)
		if err == nil {
			comment := &models.Comment{
				PostID:  uint(pid),
				Content: content,
				UserID:  userId,
			}
			err = comment.Insert()
		}
	}
	if err == nil {
		NotifyEmail("[wblog]您有一条新评论", fmt.Sprintf("<a href=\"%s/post/%d\" target=\"_blank\">%s</a>:%s", system.GetConfiguration().Domain, post.ID, post.Title, content))
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

func CommentDelete(c *gin.Context) {
	s := sessions.Default(c)
	sessionUserID := s.Get(SESSION_KEY)
	userId, _ := sessionUserID.(uint)

	commentId := c.Param("id")
	cid, err := strconv.ParseUint(commentId, 10, 64)
	if err == nil {
		comment := &models.Comment{
			UserID: uint(userId),
		}
		comment.ID = uint(cid)
		err = comment.Delete()
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": err == nil,
	})
}

func CommentRead(c *gin.Context) {
	id := c.Param("id")
	_id, err := strconv.ParseUint(id, 10, 64)
	if err == nil {
		comment := new(models.Comment)
		comment.ID = uint(_id)
		err = comment.Update()
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": err == nil,
	})
}

func CommentReadAll(c *gin.Context) {
	err := models.SetAllCommentRead()
	c.JSON(http.StatusOK, gin.H{
		"succeed": err == nil,
	})
}
