package controllers

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
)

func CommentPost(c *gin.Context) {
	var (
		err  error
		res  = gin.H{}
		post *models.Post
		cfg  = system.GetConfiguration()
	)
	defer writeJSON(c, res)
	s := sessions.Default(c)
	sessionUserID := s.Get(SessionKey)
	userId, _ := sessionUserID.(uint)

	verifyCode := c.PostForm("verifyCode")
	captchaId := s.Get(SessionCaptcha).(string)
	s.Delete(SessionCaptcha)
	if !captcha.VerifyString(captchaId, verifyCode) {
		res["message"] = "error verifyCode"
		return
	}

	content := c.PostForm("content")
	if len(content) == 0 {
		res["message"] = "content cannot be empty."
		return
	}
	pid, err := PostFormUint(c, "postId")
	if err != nil {
		res["message"] = err.Error()
		return
	}
	post, err = models.GetPostById(pid)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	comment := &models.Comment{
		PostID:  pid,
		Content: content,
		UserID:  userId,
	}
	err = comment.Insert()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	NotifyEmail(fmt.Sprintf("[%s]您有一条新评论", cfg.Title), fmt.Sprintf("<a href=\"%s/post/%d\" target=\"_blank\">%s</a>:%s", cfg.Domain, post.ID, post.Title, content))
	res["succeed"] = true
}

func CommentDelete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
		cid uint
	)
	defer writeJSON(c, res)

	s := sessions.Default(c)
	sessionUserID := s.Get(SessionKey)
	userId, _ := sessionUserID.(uint)

	cid, err = ParamUint(c, "id")
	if err != nil {
		res["message"] = err.Error()
		return
	}
	comment := &models.Comment{
		UserID: userId,
	}
	comment.ID = cid
	err = comment.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func CommentRead(c *gin.Context) {
	var (
		id  uint
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id, err = ParamUint(c, "id")
	if err != nil {
		res["message"] = err.Error()
		return
	}
	comment := new(models.Comment)
	comment.ID = id
	err = comment.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func CommentReadAll(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	err = models.SetAllCommentRead()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
