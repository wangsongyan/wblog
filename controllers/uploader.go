package controllers

import "mime/multipart"

type Uploader interface {
	upload(file multipart.File) (string, error)
}
