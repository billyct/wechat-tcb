package wechat_tcb

import "time"

type Cache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{}, time.Duration)
}