package tcb

import (
	"fmt"
	"os"
)

type Tcb struct {
	appId     string
	appSecret string
	envId     string
	cache     Cache
}

// 初始化
func New(c *Config) *Tcb {
	t := &Tcb{}
	copyConfig(t, c)
	return t
}

// 复制 config 配置到 tcb
func copyConfig(t *Tcb, c *Config) {
	if c.AppID == "" {
		c.AppID = os.Getenv("APP_ID")
	}

	if c.AppSecret == "" {
		c.AppSecret = os.Getenv("APP_SECRET")
	}

	if c.EnvID == "" {
		c.EnvID = os.Getenv("ENV_ID")
	}

	if c.Cache == nil {
		// memory cache
	}

	t.appId = c.AppID
	t.appSecret = c.AppSecret
	t.envId = c.EnvID
	t.cache = c.Cache
}

// 拼装 API URL
func (t *Tcb) url(baseUrl string) string {
	return fmt.Sprintf("%s?access_token=%s", baseUrl, t.AccessToken())
}
