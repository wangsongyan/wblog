package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
	"net/http"
	"strings"
	"time"
)

func SubscribeGet(c *gin.Context) {
	count, _ := models.CountSubscriber()
	c.HTML(http.StatusOK, "other/subscribe.html", gin.H{
		"total": count,
	})
}

func Subscribe(c *gin.Context) {
	mail := c.PostForm("mail")
	var err error
	if len(mail) > 0 {
		var subscriber *models.Subscriber
		subscriber, err = models.GetSubscriberByEmail(mail)
		if err == nil {
			if !subscriber.VerifyState && time.Now().After(subscriber.OutTime) { //激活链接超时
				err = sendActiveEmail(subscriber)
				if err == nil {
					count, _ := models.CountSubscriber()
					c.HTML(http.StatusOK, "other/subscribe.html", gin.H{
						"message": "subscribe succeed.",
						"total":   count,
					})
					return
				}
			} else if subscriber.VerifyState && !subscriber.SubscribeState { //已认证，未订阅
				subscriber.SubscribeState = true
				err = subscriber.Update()
				if err == nil {
					err = errors.New("subscribe succeed.")
				}
			} else {
				err = errors.New("mail have already actived or have unactive mail in your mailbox.")
			}
		} else {
			subscriber := &models.Subscriber{
				Email: mail,
			}
			err = subscriber.Insert()
			if err == nil {
				err = sendActiveEmail(subscriber)
				if err == nil {
					count, _ := models.CountSubscriber()
					c.HTML(http.StatusOK, "other/subscribe.html", gin.H{
						"message": "subscribe succeed.",
						"total":   count,
					})
					return
				}
			}
		}
	} else {
		err = errors.New("empty mail address.")
	}
	count, _ := models.CountSubscriber()
	c.HTML(http.StatusOK, "other/subscribe.html", gin.H{
		"message": err.Error(),
		"total":   count,
	})
}

func sendActiveEmail(subscriber *models.Subscriber) error {
	uuid := helpers.UUID()
	duration, _ := time.ParseDuration("30m")
	subscriber.OutTime = time.Now().Add(duration)
	subscriber.SecretKey = uuid
	signature := helpers.Md5(subscriber.Email + uuid + subscriber.OutTime.Format("20060102150405"))
	subscriber.Signature = signature
	err := sendMail(subscriber.Email, "[Wblog]邮箱验证", fmt.Sprintf("%s/active?sid=%s", system.GetConfiguration().Domain, signature))
	if err == nil {
		err = subscriber.Update()
	}
	return err
}

func ActiveSubsciber(c *gin.Context) {
	sid := c.Query("sid")
	var err error
	if len(sid) > 0 {
		var subscriber *models.Subscriber
		subscriber, err = models.GetSubscriberBySignature(sid)
		if err == nil {
			if time.Now().Before(subscriber.OutTime) {
				subscriber.VerifyState = true
				subscriber.OutTime = time.Now()
				err = subscriber.Update()
				if err == nil {
					HandleMessage(c, "激活成功！")
					return
				}
			} else {
				err = errors.New("激活链接已过期，请重新获取！")
			}
		} else {
			err = errors.New("激活链接有误，请重新获取！")
		}
	} else {
		err = errors.New("激活链接有误，请重新获取！")
	}
	HandleMessage(c, err.Error())
}

func UnSubscribe(c *gin.Context) {
	sid := c.Query("sid")
	if len(sid) > 0 {
		subscriber, err := models.GetSubscriberBySignature(sid)
		if err == nil && subscriber.VerifyState && subscriber.SubscribeState {
			subscriber.SubscribeState = false
			err = subscriber.Update()
			if err == nil {
				HandleMessage(c, "Unscribe Succeessful!")
				return
			}
		}
	}
	HandleMessage(c, "Internal Server Error!")
}

func GetUnSubcribeUrl(subscriber *models.Subscriber) (string, error) {
	uuid := helpers.UUID()
	signature := helpers.Md5(subscriber.Email + uuid)
	subscriber.SecretKey = uuid
	subscriber.Signature = signature
	err := subscriber.Update()
	return fmt.Sprintf("%s/unsubscribe?sid=%s", system.GetConfiguration().Domain, signature), err
}

func sendEmailToSubscribers(subject, body string) error {
	subscribers, err := models.ListSubscriber(true)
	if err == nil {
		emails := make([]string, 0)
		for _, subscriber := range subscribers {
			emails = append(emails, subscriber.Email)
		}
		if len(emails) > 0 {
			err = sendMail(strings.Join(emails, ";"), subject, body)
		} else {
			err = errors.New("no subscribers!")
		}
	}
	return err
}

func SubscriberIndex(c *gin.Context) {
	subscribers, _ := models.ListSubscriber(false)
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/subscriber.html", gin.H{
		"subscribers": subscribers,
		"user":        user,
		"comments":    models.MustListUnreadComment(),
	})
}

// 邮箱为空时，发送给所有订阅者
func SubscriberPost(c *gin.Context) {
	mail := c.PostForm("mail")
	subject := c.PostForm("subject")
	body := c.PostForm("body")
	var err error
	if len(mail) > 0 {
		err = sendMail(mail, subject, body)
	} else {
		err = sendEmailToSubscribers(subject, body)
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
