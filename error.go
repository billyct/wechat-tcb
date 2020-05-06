package tcb

import "fmt"

// API 返回的数据中的错误信息
type ResError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// API Error
type APIError struct {
	ResError
	APIName string
}

// Error
func (t *APIError) Error() string {
	return fmt.Sprintf("api error: apiName=%v errcode=%v, errmsg=%v", t.APIName, t.ErrCode, t.ErrMsg)
}

// HTTP Error
type HTTPError struct {
	Method     string
	URL        string
	StatusCode int
}

// Error
func (h *HTTPError) Error() string {
	return fmt.Sprintf("http error: method=%v url=%v, statusCode=%v", h.Method, h.URL, h.StatusCode)
}
