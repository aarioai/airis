package aconfig

import (
	"fmt"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/arrmap"
	"gopkg.in/ini.v1"
	"strings"
	"time"
)

// convertIniToMap 将 ini.FilePath 转换为扁平化的 map[string]string
func (c *Config) convertIniToMap(iniFile *ini.File, target map[string]string) error {
	for _, section := range iniFile.Sections() {
		prefix := ""
		sectionName := section.Name()
		if sectionName != "" && sectionName != "DEFAULT" {
			prefix = sectionName + "."
		}
		for _, key := range section.Keys() {
			k := prefix + key.Name()
			v := key.Value()
			if c.valueProcessor != nil {
				var err error
				if v, err = c.valueProcessor(k, v); err != nil {
					return err
				}
			}
			target[k] = v
		}
	}
	return nil
}

// loadIni 加载 ini 文件
func (c *Config) loadIni() (map[string]string, error) {
	iniFile, err := ini.Load(c.path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config %s: %w", c.path, err)
	}
	data := make(map[string]string)
	err = c.convertIniToMap(iniFile, data)
	return data, err
}
func (c *Config) processData(data map[string]string) error {
	if c.valueProcessor == nil {
		return nil
	}

	for k, v := range data {
		newV, err := c.valueProcessor(k, v)
		if err != nil {
			return err
		}
		if newV != v {
			data[k] = newV
		}
	}
	return nil
}
func (c *Config) parseOtherConfig(otherConfigs ...map[string]string) (map[string]string, error) {
	cfgs := arrmap.Merge(otherConfigs...)
	if c.valueProcessor == nil {
		return cfgs, nil
	}
	for k, v := range cfgs {
		newV, err := c.valueProcessor(k, v)
		if err != nil {
			return nil, err
		}
		if newV != v {
			cfgs[k] = newV
		}
	}
	return cfgs, nil
}

// Reload 重新加载base和text配置
func (c *Config) Reload() error {
	c.startWrite()
	defer c.endWrite()

	data, err := c.loadIni()
	if err != nil {
		return err
	}
	textConfigs, err := c.loadAllTextConfigs(data[CkTextConfigDirs])
	if err != nil {
		return err
	}

	// 写锁范围一定要越小越好
	cfgMtx.Lock()
	c.baseConfig = data
	c.textConfig = textConfigs
	cfgMtx.Unlock()
	return c.initializeConfig()
}

func (c *Config) getBase(key string) string {
	if c.isOnWrite() {
		cfgMtx.RLock()
		defer cfgMtx.RUnlock()
	}
	if key == "" || len(c.baseConfig) == 0 {
		return ""
	}

	keys := splitDots(key)
	keyName := keys[0]
	if len(keys) > 1 {
		keyName += "." + strings.Join(keys[1:], "_")
	}

	return c.baseConfig[keyName]
}

// 不要获取太细分，否则容易导致错误不容易被排查
func (c *Config) getOtherConfig(key string) string {
	if c.isOnWrite() {
		cfgMtx.RLock()
		defer cfgMtx.RUnlock()
	}
	if key == "" || len(c.otherConfig) == 0 {
		return ""
	}
	return c.otherConfig[key]
}

func (c *Config) MustGetString(key string) (string, error) {
	// 1. 优先从其他配置（如数据库下载来的）读取。能修改掉 ini 里面的模式配置
	if v := c.getOtherConfig(key); v != "" {
		return v, nil
	}
	// 2. 从 ini 配置获取
	if v := c.getBase(key); v != "" {
		return v, nil
	}
	// 3. 从text 配置读取
	if rsa := c.getTextData(key); len(rsa) > 0 {
		return rsa, nil
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

// Get value from config by key name
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

// MustGetDuration parses a valid duration string from config
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// A duration string is a possibly signed sequence of decimal numbers,
// each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
func (c *Config) MustGetDuration(key string) (time.Duration, error) {
	s, err := c.MustGetString(key)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(s)
}

// GetDuration parses a duration string from config
func (c *Config) GetDuration(key string, defaultValue time.Duration) time.Duration {
	s := c.GetString(key)
	if s == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultValue
	}
	return d
}
