package controllers

import (
	"bytes"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/system"
	"io/ioutil"
	"net/http"
	"os"
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

func RestorePost(c *gin.Context) {
	fileName := c.PostForm("fileName")
	var err error
	if len(fileName) > 0 {
		fileUrl := system.GetConfiguration().QiniuFileServer + fileName
		var resp *http.Response
		resp, err = http.Get(fileUrl)
		if err == nil {
			defer resp.Body.Close()
			var data []byte
			data, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				data, err = helpers.Decrypt(data, system.GetConfiguration().BackupKey)
				if err == nil {
					err = ioutil.WriteFile(fileName, data, os.ModePerm)
				}
			}
		}
	} else {
		err = errors.New("fileName cannot be empty.")
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

func Backup() error {
	var err error
	if exists, _ := helpers.PathExists("wblog.db"); exists {
		seelog.Debug("start backup...")

		data, err := ioutil.ReadFile("wblog.db")
		if err != nil {
			seelog.Error(err)
			return err
		}
		encryptData, err := helpers.Encrypt(data, system.GetConfiguration().BackupKey)
		if err != nil {
			seelog.Error(err)
			return err
		}

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
		err = uploader.Put(nil, &ret, token, fileName, bytes.NewReader(encryptData), int64(len(encryptData)), nil)
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
