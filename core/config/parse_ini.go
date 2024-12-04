package config

import (
	"fmt"
	"github.com/aarioai/airis/core/atype"
	"gopkg.in/ini.v1"
	"strings"
	"sync"
	"time"
)

var (
	cfgMtx sync.RWMutex
)

func New(path string) *Config {
	cfg := &Config{path: path}
	cfg.Reload()
	return cfg
}

func (c *Config) Reload() {
	data, err := ini.Load(c.path)
	if err != nil {
		panic(err)
	}
	cfgMtx.Lock()
	defer cfgMtx.Unlock()

	c.data = data

	c.Env = Env(c.GetString(CkEnv))
	c.TimezoneID, _ = time.Now().Zone()
	c.TimeLocation = time.Local
	if tz := c.GetString(CkTimezoneID); tz != "" {
		c.TimezoneID = tz
		c.TimeLocation, err = time.LoadLocation(tz)
		if err != nil {
			panic("invalid timezone: " + tz + ", error: " + err.Error())
		}
	}
	c.TimeFormat = c.GetString(CkTimeFormat, "2006-02-01 15:04:05")
	c.Mock, _ = c.Get(CkMock).Bool()

	if err = c.loadRsa(); err != nil {
		panic(err)
	}
}

// 这里有锁，所以要批量设置
func (c *Config) AddConfigs(otherConfigs map[string]string) {
	cfgMtx.Lock()
	defer cfgMtx.Unlock()
	if otherConfigs == nil {
		c.otherConfig = make(map[string]string)
	}
	for k, v := range otherConfigs {
		c.otherConfig[k] = v
	}
}

func (c *Config) getIni(key string) string {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()

	keys := splitDots(key)
	var s *ini.Section
	if len(keys) == 1 {
		if s = c.data.Section(""); s.HasKey(key) {
			return s.Key(key).String()
		}
		return ""
	}

	k := strings.Join(keys[1:], "_")
	if s = c.data.Section(keys[0]); s.HasKey(k) {
		return s.Key(k).String()
	}
	return ""
}

// 不要获取太细分，否则容易导致错误不容易被排查
func (c *Config) getOtherConfig(key string) string {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()
	d, _ := c.otherConfig[key]
	return d
}

func (c *Config) MustGetString(key string) (string, error) {
	v := c.getIni(key)
	if v != "" {
		return v, nil
	}
	// 从RSA读取
	if rsa, _ := c.getRsa(key); len(rsa) > 0 {
		return string(rsa), nil
	}
	// 从其他配置（如数据库下载来的）读取
	if v = c.getOtherConfig(key); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("must set config `%s`", key)
}

func (c *Config) GetString(key string, defaultValue ...string) string {
	v, _ := c.MustGetString(key)
	if v != "" {
		return v
	}
	if len(defaultValue) > 0 {
		v = defaultValue[0]
	}
	return v
}

func (c *Config) MustGet(key string) (*atype.Atype, error) {
	v, err := c.MustGetString(key)
	if err != nil {
		return nil, err
	}
	return atype.New(v), nil
}

// Get(key) or Get(key, defaultValue)
// 先从 ini 文件读取，找不到再去从其他 provider （如数据库拉下来的配置）里面找
func (c *Config) Get(key string, defaultValue ...any) *atype.Atype {
	v, _ := c.MustGetString(key)
	if v != "" {
		return atype.New(v)
	}
	if len(defaultValue) > 0 {
		return atype.New(defaultValue[0])
	}
	return atype.New("")
}
