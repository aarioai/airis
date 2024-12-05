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
func (c *Config) reload() {
	data, err := ini.Load(c.path)
	if err != nil {
		panic(fmt.Errorf("failed to load config file %s: %w", c.path, err))
	}
	// 写锁范围一定要越小越好
	cfgMtx.Lock()
	c.data = data
	cfgMtx.Unlock()
}

func (c *Config) initializeConfig() error {
	c.Env = Env(c.GetString(CkEnv))
	c.Mock, _ = c.Get(CkMock).Bool()
	// 初始化时区配置
	if err := c.initializeTimezone(); err != nil {
		return err
	}
	return c.loadRsa()
}
func (c *Config) initializeTimezone() error {
	var err error
	c.TimezoneID, _ = time.Now().Zone()
	c.TimeLocation = time.Local
	c.TimeFormat = c.GetString(CkTimeFormat, "2006-02-01 15:04:05")

	if tz := c.GetString(CkTimezoneID); tz != "" {
		c.TimezoneID = tz
		if c.TimeLocation, err = time.LoadLocation(tz); err != nil {
			return fmt.Errorf("invalid timezone %s: %w", tz, err)
		}
	}
	return nil
}

func (c *Config) Reload() {
	c.reload()
	if err := c.initializeConfig(); err != nil {
		panic(fmt.Errorf("failed to initialize config: %w", err))
	}
}

// 这里有锁，所以要批量设置
func (c *Config) AddConfigs(otherConfigs map[string]string) {
	cfgMtx.Lock()
	defer cfgMtx.Unlock()
	if c.otherConfig == nil {
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
	if len(keys) == 1 {
		return c.data.Section("").Key(key).String()
	}
	section := c.data.Section(keys[0])
	return section.Key(strings.Join(keys[1:], "_")).String()
}

// 不要获取太细分，否则容易导致错误不容易被排查
func (c *Config) getOtherConfig(key string) string {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()
	return c.otherConfig[key]
}

func (c *Config) MustGetString(key string) (string, error) {
	if v := c.getIni(key); v != "" {
		return v, nil
	}
	// 从RSA读取
	if rsa, _ := c.getRsa(key); len(rsa) > 0 {
		return string(rsa), nil
	}
	// 从其他配置（如数据库下载来的）读取
	if v := c.getOtherConfig(key); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("required config key not found: %s", key)
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
