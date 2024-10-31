package controllers

import (
	"context"
	"mime/multipart"
	"os"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/wangsongyan/wblog/system"
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

func (u QiniuUploader) upload(file multipart.File, _ *multipart.FileHeader) (url string, err error) {
	var (
		ret  PutRet
		size int64
		cfg  = system.GetConfiguration()
	)
	if statInterface, ok := file.(Stat); ok {
		fileInfo, _ := statInterface.Stat()
		size = fileInfo.Size()
	}
	if sizeInterface, ok := file.(Size); ok {
		size = sizeInterface.Size()
	}

	putPolicy := storage.PutPolicy{
		Scope: cfg.Qiniu.Bucket,
	}
	mac := qbox.NewMac(cfg.Qiniu.AccessKey, cfg.Qiniu.SecretKey)
	token := putPolicy.UploadToken(mac)
	uploader := storage.NewFormUploader(&storage.Config{})
	putExtra := storage.PutExtra{}

	err = uploader.PutWithoutKey(context.Background(), &ret, token, file, size, &putExtra)
	if err != nil {
		return
	}
	if strings.HasSuffix(cfg.Qiniu.FileServer, "/") {
		url = cfg.Qiniu.FileServer + ret.Key
	} else {
		url = cfg.Qiniu.FileServer + "/" + ret.Key
	}
	return
}
