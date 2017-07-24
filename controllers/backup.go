package controllers

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/system"
	"net/http"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
	"time"
)

func BackupPost(c *gin.Context) {
	err := Backup()
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

func Backup() error {
	var err error
	if exists, _ := helpers.PathExists("wblog.db"); exists {
		seelog.Debug("start backup...")

		conf.ACCESS_KEY = system.GetConfiguration().QiniuAccessKey
		conf.SECRET_KEY = system.GetConfiguration().QiniuSecretKey

		// 创建一个Client
		c := kodo.New(0, nil)
		// 设置上传的策略
		policy := &kodo.PutPolicy{
			Scope: system.GetConfiguration().QiniuBucket,
			//设置Token过期时间
			Expires: 3600,
		}
		// 生成一个上传token
		token := c.MakeUptoken(policy)
		// 构建一个uploader
		zone := 0
		uploader := kodocli.NewUploader(zone, nil)

		var ret PutRet
		fileName := fmt.Sprintf("wblog_%s.db", time.Now().Format("20060102150405"))
		err = uploader.PutFile(nil, &ret, token, fileName, "wblog.db", nil)
		if err == nil {
			seelog.Debug("backup succeefully.")
		} else {
			seelog.Debugf("backup error:%v", err)
		}
	} else {
		err = errors.New("database file doesn't exists.")
		seelog.Debug("database file doesn't exists.")
	}
	return err
}
