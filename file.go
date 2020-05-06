package tcb

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

const (
	// 请求地址：「获取文件上传链接」
	urlUploadFile = "https://api.weixin.qq.com/tcb/uploadfile"
)

// 获取文件上传链接请求
type reqUploadFile struct {
	Env  string `json:"env,omitempty"`
	Path string `json:"path,omitempty"`
}

// 获取文件上传链接返回
type resUploadFile struct {
	ResError
	URL           string `json:"url"`           //上传url
	Token         string `json:"token"`         //token
	Authorization string `json:"authorization"` //authorization
	FileID        string `json:"file_id"`       //文件ID
	CosFileID     string `json:"cos_file_id"`   //cos文件ID
}

// 获取文件上传链接，并上传文件
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/storage/uploadFile.html
func (t *Tcb) UploadFileWithFile(path string, file string) (string, error) {
	res, err := t.UploadFile(path)
	if err != nil {
		return "", err
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	// 增加其他属性
	w.WriteField("key", path)
	w.WriteField("Signature", res.Authorization)
	w.WriteField("x-cos-security-token", res.Token)
	w.WriteField("x-cos-meta-fileid", res.CosFileID)

	// 增加 file 二进制文件内容
	part, err := w.CreateFormFile("file", file)
	if err != nil {
		return "", err
	}

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, f)
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	_, err = httpPost(res.URL, w.FormDataContentType(), body)
	if err != nil {
		return "", err
	}

	return res.FileID, nil
}

// 获取文件上传链接
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/storage/uploadFile.html
func (t *Tcb) UploadFile(path string) (*resUploadFile, error) {
	req := &reqUploadFile{
		Env:  t.envId,
		Path: path,
	}

	data, err := httpPostWithReq(t.url(urlUploadFile), req)
	if err != nil {
		return nil, err
	}

	res := &resUploadFile{}

	err = DecodeApiData("UploadFile", data, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
