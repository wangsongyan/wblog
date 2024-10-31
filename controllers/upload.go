package controllers

import (
	"github.com/wangsongyan/wblog/system"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	var (
		err      error
		res      = gin.H{}
		url      string
		uploader Uploader
		file     multipart.File
		fh       *multipart.FileHeader
		cfg      = system.GetConfiguration()
	)
	defer writeJSON(c, res)
	file, fh, err = c.Request.FormFile("file")
	if err != nil {
		res["message"] = err.Error()
		return
	}

	if cfg.FileServer == "smms" && cfg.Smms.Enabled {
		uploader = SmmsUploader{}
	}
	if cfg.FileServer == "qiniu" && cfg.Qiniu.Enabled {
		uploader = QiniuUploader{}
	}
	url, err = uploader.upload(file, fh)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
	res["url"] = url
}
