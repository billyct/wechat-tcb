package tcb

import (
	"fmt"
	"time"
)

const (
	// 请求地址：「接口调用凭证」
	urlGetAccessToken = "https://api.weixin.qq.com/cgi-bin/token"
	// 缓存 access_token 的 key
	cacheKeyAccessToken = "wechat-tcb:access_token"
)

// 接口调用凭证返回
type resAccessToken struct {
	ResError
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// 返回 access token，无错误
func (t *Tcb) AccessToken() string {
	token, _ := t.GetAccessToken()
	return token
}

// 接口调用凭证，返回 access_token，并缓存
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html
func (t *Tcb) GetAccessToken() (string, error) {

	// 如果有缓存则取缓存中的 token
	if v, ok := t.cache.Get(cacheKeyAccessToken); ok {
		return v.(string), nil
	}

	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", urlGetAccessToken, t.appId, t.appSecret)

	data, err := httpGet(url)
	if err != nil {
		return "", err
	}

	var res resAccessToken

	err = DecodeApiData("GetAccessToken", data, &res)
	if err != nil {
		return "", err
	}

	t.cache.Set(cacheKeyAccessToken, res.AccessToken, time.Duration(res.ExpiresIn)*time.Second)

	return res.AccessToken, nil
}
