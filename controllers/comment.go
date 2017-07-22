package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/models"
	"net/http"
	"strconv"
)

func CommentPost(c *gin.Context) {
	s := sessions.Default(c)
	sessionUserID := s.Get(SESSION_KEY)
	userId, _ := sessionUserID.(uint)

	postId := c.PostForm("postId")
	content := c.PostForm("content")
	pid, err := strconv.ParseUint(postId, 10, 64)
	if err == nil {
		comment := &models.Comment{
			PostID:  uint(pid),
			Content: content,
			UserID:  userId,
		}
		comment.Insert()
	}
	c.Redirect(http.StatusMovedPermanently, "/post/"+postId)
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
