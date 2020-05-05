package tcb

import "fmt"

type Error struct {
	Api     string
	ErrCode int64
	ErrMsg  string
}

func (t *Error) Error() string {
	return fmt.Sprintf("api % error, errcode: %s, errmsg: %s", t.Api, t.ErrCode, t.ErrMsg)
}
