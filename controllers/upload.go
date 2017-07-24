package controllers

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/wangsongyan/wblog/helpers"
	"github.com/wangsongyan/wblog/system"
	"mime/multipart"
	"net/http"
	"os"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
	"time"
)

// 获取文件大小的接口
type Size interface {
	Size() int64
}

// 获取文件信息的接口
type Stat interface {
	Stat() (os.FileInfo, error)
}

// 构造返回值字段
type PutRet struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

func Upload(c *gin.Context) {
	var err error
	file, _, err := c.Request.FormFile("file")
	if err == nil {
		var key string
		key, err = uploadFile(file)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"succeed": true,
				"url":     system.GetConfiguration().QiniuFileServer + key,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": false,
		"message": err.Error(),
	})
}

func uploadFile(file multipart.File) (string, error) {

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

	var size int64
	if statInterface, ok := file.(Stat); ok {
		fileInfo, _ := statInterface.Stat()
		size = fileInfo.Size()
	}
	if sizeInterface, ok := file.(Size); ok {
		size = sizeInterface.Size()
	}

	var ret PutRet
	err := uploader.PutWithoutKey(nil, &ret, token, file, size, nil)
	if err != nil {
		return "", err
	}
	return ret.Key, nil
}

func Backup() {
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
		err := uploader.PutFile(nil, &ret, token, fileName, "wblog.db", nil)
		if err == nil {
			seelog.Debug("backup succeefully.")
		} else {
			seelog.Debugf("backup error:%v", err)
		}
	} else {
		seelog.Debug("database file doesn't exists.")
	}
}
