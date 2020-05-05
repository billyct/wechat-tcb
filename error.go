package wechat_tcb

import "fmt"

type TcbError struct {
	Api     string
	ErrCode int64
	ErrMsg  string
}

func (t *TcbError) Error() string {
	return fmt.Sprintf("api % error, errcode: %s, errmsg: %s", t.Api, t.ErrCode, t.ErrMsg)
}
