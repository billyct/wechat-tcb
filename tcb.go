package tcb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	kAccessToken = "wechat-tcb:access_token"
)

type tcbOption func(t *Tcb)

// 设置 cache 属性
func WithCache(c Cache) tcbOption {
	return func(t *Tcb) {
		t.cache = c
	}
}

type Tcb struct {
	appId     string
	appSecret string
	cache     Cache
}

// 初始化
func New(options ...tcbOption) *Tcb {
	t := &Tcb{
		appId:     os.Getenv("APP_ID"),
		appSecret: os.Getenv("APP_SECRET"),
	}

	for _, option := range options {
		option(t)
	}

	return t
}

// 请求地址：「接口调用凭证」
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html
func (t *Tcb) urlGetAccessToken() string {
	return fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", t.appId, t.appSecret)
}

// 请求地址：「获取文件上传链接」
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/storage/uploadFile.html
func (t *Tcb) urlUploadFile() string {
	return fmt.Sprintf("https://api.weixin.qq.com/tcb/uploadfile?access_token=%s", t.AccessToken())
}

// 请求地址：「新增集合」
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/database/databaseCollectionAdd.html
func (t *Tcb) urlDatabaseCollectionAdd() string {
	return fmt.Sprintf("https://api.weixin.qq.com/tcb/databasecollectionadd?access_token=%s", t.AccessToken())
}

// 请求地址：「数据库插入记录」
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/database/databaseAdd.html
func (t *Tcb) urlDataBaseAdd() string {
	return fmt.Sprintf("https://api.weixin.qq.com/tcb/databaseadd?access_token=%s", t.AccessToken())
}

// 发送 GET 请求
func (t *Tcb) get(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return t.parseResponse(res)
}

// 发送 POST 请求
func (t *Tcb) post(url string, data []byte) (map[string]interface{}, error) {
	res, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return t.parseResponse(res)
}

// 解析请求返回的数据为 map[string]interface{}
func (t *Tcb) parseResponse(res *http.Response) (map[string]interface{}, error) {
	var ret map[string]interface{}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(data, &ret)
	return ret, err
}

// 检查请求微信返回的错误
func (t *Tcb) checkError(name string, res map[string]interface{}) (error, bool) {
	if code, ok := res["errcode"]; ok {
		errCode := int64(code.(float64))

		if errCode != 0 {
			err := &Error{
				Api:     name,
				ErrCode: errCode,
				ErrMsg:  res["errmsg"].(string),
			}

			return err, true
		}
	}

	return nil, false
}

// 返回 access token，无错误
func (t *Tcb) AccessToken() string {
	token, _ := t.GetAccessToken()
	return token
}

// 获取「云环境ID」
func (t *Tcb) EnvId() string {
	return os.Getenv("ENV_ID")
}

// 封装需要「云环境ID」的 POST 请求数据
func (t *Tcb) postDataWithEnvId(data map[string]interface{}) ([]byte, error) {
	ret := map[string]interface{}{
		"env": t.EnvId(),
	}

	for k, v := range data {
		ret[k] = v
	}

	return json.Marshal(ret)
}

// 接口调用凭证，返回 access_token，并缓存
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html
func (t *Tcb) GetAccessToken() (string, error) {

	// 如果有缓存则取缓存中的 token
	if v, ok := t.cache.Get(kAccessToken); ok {
		return v.(string), nil
	}

	res, err := t.get(t.urlGetAccessToken())
	if err != nil {
		return "", err
	}

	accessToken := res["access_token"].(string)
	expiresIn := res["expires_in"].(float64)
	// 缓存 res
	t.cache.Set(kAccessToken, accessToken, time.Duration(expiresIn)*time.Second)

	return accessToken, nil
}

// 上传链接使用
// 拼装一个上传文件的 HTTP POST 请求
func (t *Tcb) requestUploadFile(res map[string]interface{}, key string, file *os.File) (*http.Request, error) {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	// 增加其他属性
	writer.WriteField("key", key)
	writer.WriteField("Signature", res["authorization"].(string))
	writer.WriteField("x-cos-security-token", res["token"].(string))
	writer.WriteField("x-cos-meta-fileid", res["cos_file_id"].(string))

	// 增加 file 二进制文件内容
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	uri := res["url"].(string)
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

// 获取文件上传链接，并上传文件
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/storage/uploadFile.html
func (t *Tcb) UploadFile(key, file string) (string, error) {
	// 请求参数
	data, err := t.postDataWithEnvId(map[string]interface{}{
		"path": key,
	})
	if err != nil {
		return "", err
	}

	res, err := t.post(t.urlUploadFile(), data)
	if err != nil {
		return "", err
	}

	if err, yes := t.checkError("UploadFile", res); yes {
		return "", err
	}

	// 获取上传链接请求成功
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}

	defer f.Close()

	req, err := t.requestUploadFile(res, key, f)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return "", err
	}

	return res["file_id"].(string), nil
}

// 新增集合
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/database/databaseCollectionAdd.html
func (t *Tcb) DatabaseCollectionAdd(name string) error {
	data, err := t.postDataWithEnvId(map[string]interface{}{
		"collection_name": name,
	})
	if err != nil {
		return err
	}

	res, err := t.post(t.urlDatabaseCollectionAdd(), data)
	if err != nil {
		return err
	}

	if err, yes := t.checkError("DatabaseCollectionAdd", res); yes {
		return err
	}

	return nil
}

// 数据库插入记录
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/database/databaseAdd.html
func (t *Tcb) DatabaseAdd(query string) error {
	data, err := t.postDataWithEnvId(map[string]interface{}{
		"query": query,
	})
	if err != nil {
		return err
	}

	res, err := t.post(t.urlDataBaseAdd(), data)
	if err != nil {
		return err
	}

	if err, yes := t.checkError("DatabaseAdd", res); yes {
		return err
	}

	return nil
}