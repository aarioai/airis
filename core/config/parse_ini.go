package config

import (
	"fmt"
	"github.com/aarioai/airis/core/atype"
	"github.com/aarioai/airis/pkg/arrmap"
	"gopkg.in/ini.v1"
	"strings"
	"sync"
)

var (
	cfgMtx sync.RWMutex
)

func New(path string, otherConfigs ...map[string]string) *Config {
	cfg := &Config{
		path: path,
	}
	if err := cfg.Reload(otherConfigs...); err != nil {
		panic(err.Error())
	}
	return cfg
}
func (c *Config) isOnWrite() bool {
	return c.onWrite.Load()
}

// convertIniToMap 将 ini.File 转换为扁平化的 map[string]string
func convertIniToMap(iniFile *ini.File, target map[string]string) {
	// 处理默认section
	defaultSection := iniFile.Section("")
	for _, key := range defaultSection.Keys() {
		target[key.Name()] = key.Value()
	}

	// 处理其他sections
	for _, section := range iniFile.Sections() {
		// 跳过默认section，因为已经处理过了
		if section.Name() == "" {
			continue
		}

		// 遍历section中的所有键值对
		for _, key := range section.Keys() {
			// 使用 section.name.key 作为完整的键名
			fullKey := section.Name() + "." + key.Name()
			target[fullKey] = key.Value()
		}
	}
}

func (c *Config) getIni(key string) string {
	if c.isOnWrite() {
		cfgMtx.RLock()
		defer cfgMtx.RUnlock()
	}
	if key == "" || len(c.data) == 0 {
		return ""
	}

	keys := splitDots(key)
	keyName := keys[0]
	if len(keys) > 1 {
		keyName += "." + strings.Join(keys[1:], "_")
	}

	return c.data[keyName]
}

func (c *Config) loadIni() (map[string]string, error) {
	iniFile, err := ini.Load(c.path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config %s: %w", c.path, err)
	}
	data := make(map[string]string)
	convertIniToMap(iniFile, data)
	return data, nil
}

func (c *Config) Reload(otherConfigs ...map[string]string) error {
	c.startWrite()
	defer c.endWrite()

	data, err := c.loadIni()
	if err != nil {
		return err
	}
	rsa, err := c.loadRsa()
	if err != nil {
		return err
	}
	err = c.initializeConfig()
	if err != nil {
		return err
	}
	// 写锁范围一定要越小越好
	cfgMtx.Lock()
	c.data = data
	c.rsa = rsa
	// clear(c.otherConfig)
	c.otherConfig = arrmap.Merge(otherConfigs...)
	cfgMtx.Unlock()
	return nil
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
	if v := c.getIni(key); v != "" {
		return v, nil
	}
	// 从RSA读取
	if rsa := c.getRsa(key); len(rsa) > 0 {
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
