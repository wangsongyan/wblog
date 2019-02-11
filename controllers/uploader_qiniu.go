package controllers

import (
	"mime/multipart"
	"os"

	"github.com/wangsongyan/wblog/system"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
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

type QiniuUploader struct {
}

func (u QiniuUploader) upload(file multipart.File) (url string, err error) {

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
	err = uploader.PutWithoutKey(nil, &ret, token, file, size, nil)
	if err != nil {
		return
	}
	url = system.GetConfiguration().QiniuFileServer + ret.Key
	return
}
