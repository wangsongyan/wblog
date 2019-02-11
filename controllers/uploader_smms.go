package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/wangsongyan/wblog/models"
	"github.com/wangsongyan/wblog/system"
)

type SmmsUploader struct {
}

func (u SmmsUploader) upload(file multipart.File) (url string, err error) {
	var (
		resp      *http.Response
		bodyBytes []byte
		ret       models.SmmsFile
	)
	resp, err = http.Post(system.GetConfiguration().SmmsFileServer, "multipart/form-data", file)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &ret)
	if err != nil {
		return
	}
	if ret.Code == "error" {
		err = errors.New(ret.Msg)
		return
	}
	err = ret.Insert()
	if err != nil {
		return
	}
	url = ret.Url
	return
}
