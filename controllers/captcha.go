package controllers

import (
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CaptchaGet(context *gin.Context) {
	session := sessions.Default(context)
	captchaId := captcha.NewLen(4)
	session.Delete(SessionCaptcha)
	session.Set(SessionCaptcha, captchaId)
	session.Save()
	captcha.WriteImage(context.Writer, captchaId, 100, 40)
}
