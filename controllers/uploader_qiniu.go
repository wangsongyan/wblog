package controllers

import (
	"mime/multipart"
	"os"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
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

func (u QiniuUploader) upload(file multipart.File, fileHeader *multipart.FileHeader) (url string, err error) {
	var (
		ret  PutRet
		size int64
	)
	if statInterface, ok := file.(Stat); ok {
		fileInfo, _ := statInterface.Stat()
		size = fileInfo.Size()
	}
	if sizeInterface, ok := file.(Size); ok {
		size = sizeInterface.Size()
	}

	putPolicy := storage.PutPolicy{
		Scope: system.GetConfiguration().QiniuBucket,
	}
	mac := qbox.NewMac(system.GetConfiguration().QiniuAccessKey, system.GetConfiguration().QiniuSecretKey)
	token := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	uploader := storage.NewFormUploader(&cfg)
	err = uploader.PutWithoutKey(nil, &ret, token, file, size, nil)
	if err != nil {
		return
	}
	url = system.GetConfiguration().QiniuFileServer + ret.Key
	return
}
