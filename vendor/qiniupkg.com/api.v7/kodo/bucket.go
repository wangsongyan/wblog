package kodo

import (
	"encoding/base64"
	"io"
	"net/url"
	"strconv"

	. "golang.org/x/net/context"
	"qiniupkg.com/api.v7/api"
	"qiniupkg.com/x/log.v7"
)

// ----------------------------------------------------------

// 批量操作。
//
func (p *Client) Batch(ctx Context, ret interface{}, op []string) (err error) {

	return p.CallWithForm(ctx, ret, "POST", p.RSHost+"/batch", map[string][]string{"op": op})
}

// ----------------------------------------------------------

type Bucket struct {
	api.BucketInfo
	Conn *Client
	Name string
}

// 取七牛空间（bucket）的对象实例。
//
// name 是创建该七牛空间（bucket）时采用的名称。
//
func (p *Client) Bucket(name string) Bucket {
	b, err := p.BucketWithSafe(name)
	if err != nil {
		log.Errorf("Bucket(%s) failed: %+v", name, err)
	}
	return b
}

func (p *Client) BucketWithSafe(name string) (Bucket, error) {
	var info api.BucketInfo
	if len(p.UpHosts) == 0 {
		var err error
		info, err = p.apiCli.GetBucketInfo(p.mac.AccessKey, name)
		if err != nil {
			return Bucket{}, err
		}
	} else {
		info.IoHost = p.IoHost
		info.UpHosts = p.UpHosts
	}
	return Bucket{info, p, name}, nil
}

type Entry struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	EndUser  string `json:"endUser"`
}

// 取文件属性。
//
// ctx 是请求的上下文。
// key 是要访问的文件的访问路径。
//
func (p Bucket) Stat(ctx Context, key string) (entry Entry, err error) {
	err = p.Conn.Call(ctx, &entry, "POST", p.Conn.RSHost+URIStat(p.Name, key))
	return
}

// 删除一个文件。
//
// ctx 是请求的上下文。
// key 是要删除的文件的访问路径。
//
func (p Bucket) Delete(ctx Context, key string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.Conn.RSHost+URIDelete(p.Name, key))
}

// 移动一个文件。
//
// ctx     是请求的上下文。
// keySrc  是要移动的文件的旧路径。
// keyDest 是要移动的文件的新路径。
//
func (p Bucket) Move(ctx Context, keySrc, keyDest string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.Conn.RSHost+URIMove(p.Name, keySrc, p.Name, keyDest))
}

// 跨空间（bucket）移动一个文件。
//
// ctx        是请求的上下文。
// keySrc     是要移动的文件的旧路径。
// bucketDest 是文件的目标空间。
// keyDest    是要移动的文件的新路径。
//
func (p Bucket) MoveEx(ctx Context, keySrc, bucketDest, keyDest string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.Conn.RSHost+URIMove(p.Name, keySrc, bucketDest, keyDest))
}

// 复制一个文件。
//
// ctx     是请求的上下文。
// keySrc  是要复制的文件的源路径。
// keyDest 是要复制的文件的目标路径。
//
func (p Bucket) Copy(ctx Context, keySrc, keyDest string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.Conn.RSHost+URICopy(p.Name, keySrc, p.Name, keyDest))
}

// 修改文件的MIME类型。
//
// ctx  是请求的上下文。
// key  是要修改的文件的访问路径。
// mime 是要设置的新MIME类型。
//
func (p Bucket) ChangeMime(ctx Context, key, mime string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.Conn.RSHost+URIChangeMime(p.Name, key, mime))
}

// 从网上抓取一个资源并存储到七牛空间（bucket）中。
//
// ctx 是请求的上下文。
// key 是要存储的文件的访问路径。如果文件已经存在则覆盖。
// url 是要抓取的资源的URL。
//
func (p Bucket) Fetch(ctx Context, key string, url string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.IoHost+uriFetch(p.Name, key, url))
}

// ----------------------------------------------------------

type ListItem struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	EndUser  string `json:"endUser"`
}

// 首次请求，请将 marker 设置为 ""。
// 无论 err 值如何，均应该先看 entries 是否有内容。
// 如果后续没有更多数据，err 返回 EOF，markerOut 返回 ""（但不通过该特征来判断是否结束）。
//
func (p Bucket) List(
	ctx Context, prefix, delimiter, marker string, limit int) (entries []ListItem, commonPrefixes []string, markerOut string, err error) {

	listUrl := p.makeListURL(prefix, delimiter, marker, limit)

	var listRet struct {
		Marker   string     `json:"marker"`
		Items    []ListItem `json:"items"`
		Prefixes []string   `json:"commonPrefixes"`
	}
	err = p.Conn.Call(ctx, &listRet, "POST", listUrl)
	if err != nil {
		return
	}
	if listRet.Marker == "" {
		return listRet.Items, listRet.Prefixes, "", io.EOF
	}
	return listRet.Items, listRet.Prefixes, listRet.Marker, nil
}

func (p Bucket) makeListURL(prefix, delimiter, marker string, limit int) string {

	query := make(url.Values)
	query.Add("bucket", p.Name)
	if prefix != "" {
		query.Add("prefix", prefix)
	}
	if delimiter != "" {
		query.Add("delimiter", delimiter)
	}
	if marker != "" {
		query.Add("marker", marker)
	}
	if limit > 0 {
		query.Add("limit", strconv.FormatInt(int64(limit), 10))
	}
	return p.Conn.RSFHost + "/list?" + query.Encode()
}

// ----------------------------------------------------------

type BatchStatItemRet struct {
	Data  Entry  `json:"data"`
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func (p Bucket) BatchStat(ctx Context, keys ...string) (ret []BatchStatItemRet, err error) {

	b := make([]string, len(keys))
	for i, key := range keys {
		b[i] = URIStat(p.Name, key)
	}
	err = p.Conn.Batch(ctx, &ret, b)
	return
}

type BatchItemRet struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func (p Bucket) BatchDelete(ctx Context, keys ...string) (ret []BatchItemRet, err error) {

	b := make([]string, len(keys))
	for i, key := range keys {
		b[i] = URIDelete(p.Name, key)
	}
	err = p.Conn.Batch(ctx, &ret, b)
	return
}

type KeyPair struct {
	Src  string
	Dest string
}

func (p Bucket) BatchMove(ctx Context, entries ...KeyPair) (ret []BatchItemRet, err error) {

	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URIMove(p.Name, e.Src, p.Name, e.Dest)
	}
	err = p.Conn.Batch(ctx, &ret, b)
	return
}

func (p Bucket) BatchCopy(ctx Context, entries ...KeyPair) (ret []BatchItemRet, err error) {

	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URICopy(p.Name, e.Src, p.Name, e.Dest)
	}
	err = p.Conn.Batch(ctx, &ret, b)
	return
}

// ----------------------------------------------------------

func encodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

func uriFetch(bucket, key, url string) string {
	return "/fetch/" + encodeURI(url) + "/to/" + encodeURI(bucket+":"+key)
}

func URIDelete(bucket, key string) string {
	return "/delete/" + encodeURI(bucket+":"+key)
}

func URIStat(bucket, key string) string {
	return "/stat/" + encodeURI(bucket+":"+key)
}

func URICopy(bucketSrc, keySrc, bucketDest, keyDest string) string {
	return "/copy/" + encodeURI(bucketSrc+":"+keySrc) + "/" + encodeURI(bucketDest+":"+keyDest)
}

func URIMove(bucketSrc, keySrc, bucketDest, keyDest string) string {
	return "/move/" + encodeURI(bucketSrc+":"+keySrc) + "/" + encodeURI(bucketDest+":"+keyDest)
}

func URIChangeMime(bucket, key, mime string) string {
	return "/chgm/" + encodeURI(bucket+":"+key) + "/mime/" + encodeURI(mime)
}

// ----------------------------------------------------------
