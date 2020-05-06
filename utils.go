package tcb

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

// HTTP GET 请求
func httpGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			Method:     "GET",
			URL:        url,
			StatusCode: res.StatusCode,
		}
	}

	return ioutil.ReadAll(res.Body)
}

// HTTP POST 请求
func httpPost(url, contentType string, body io.Reader) ([]byte, error)  {
	res, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			Method:     "POST",
			URL:        url,
			StatusCode: res.StatusCode,
		}
	}

	return ioutil.ReadAll(res.Body)
}

// HTTP POST 请求
// Content-Type 为 application/json
func httpPostWithReq(url string, req interface{}) ([]byte, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return httpPost(url, "application/json;charset=utf-8", bytes.NewReader(data))
}

// 解析 API 数据到 obj
// 如果 API 返回错误则返回一个 APIError
func DecodeApiData(apiName string, data []byte, obj interface{}) error {
	err := json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	responseError := reflect.ValueOf(obj).Elem().FieldByName("ResError")

	code := responseError.FieldByName("ErrCode")
	if code.Int() != 0 {
		return &APIError{
			APIName:  apiName,
			ResError: responseError.Interface().(ResError),
		}
	}

	return nil
}
