package config

import (
	"fmt"
	"github.com/aarioai/airis/pkg/arrmap"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func (c *Config) AddConfigs(otherConfigs ...map[string]string) {
	c.startWrite()
	defer c.endWrite()
	cfgs := arrmap.Merge(otherConfigs...)
	if len(cfgs) == 0 {
		return
	}

	// 写锁范围一定要越小越好
	cfgMtx.Lock()
	if c.otherConfig == nil {
		c.otherConfig = make(map[string]string)
	}
	for k, v := range cfgs {
		c.otherConfig[k] = v
	}
	cfgMtx.Unlock()
}

func (c *Config) AddBinConfigs(rsaConfigs ...map[string][]byte) {
	c.startWrite()
	defer c.endWrite()
	cfgs := arrmap.MergeSlices(rsaConfigs...)
	if len(cfgs) == 0 {
		return
	}
	// 写锁范围一定要越小越好
	cfgMtx.Lock()
	if c.binConfig == nil {
		c.binConfig = make(map[string][]byte)
	}
	for name, v := range cfgs {
		if len(v) > 0 {
			c.binConfig[name] = v
		}
	}
	cfgMtx.Unlock()
}

func (c *Config) startWrite() {
	for {
		// 更新过程中，可能有其他线程也在更新（理论上几乎不存在并发重置配置，但是这样处理并发性能消耗不大，所以就这么处理下）
		// 比较 最新的 c.onWrite.Load() 等于false。若不等于false，则表示其他写操作正在进行，需要等待。若为false，则置为true
		if c.onWrite.CompareAndSwap(false, true) {
			break
		}
	}
}
func (c *Config) endWrite() {
	for {
		if c.onWrite.CompareAndSwap(true, false) {
			break
		}
	}
}

func (c *Config) initializeConfig() error {
	c.Env = Env(c.GetString(CkEnv))
	c.Mock, _ = c.Get(CkMock).Bool()
	c.FileConfigDirs = binConfigDirs(c.GetString(CkBinConfigDirs))
	// 初始化时区配置
	return c.initializeTimezone()
}

func (c *Config) initializeTimezone() error {
	var err error

	c.TimeFormat = c.GetString(CkTimeFormat, "2006-02-01 15:04:05")
	c.TimezoneID, _ = time.Now().Zone()
	c.TimeLocation = time.Local

	if tz := c.GetString(CkTimezoneID); tz != "" {
		c.TimezoneID = tz
		if c.TimeLocation, err = time.LoadLocation(tz); err != nil {
			return fmt.Errorf("invalid timezone %s: %w", tz, err)
		}
	}
	return nil
}

// 将空格隔开的配置，转换为数组
func (c *Config) loadAllBinConfigs(value string) (map[string][]byte, error) {
	dirs := binConfigDirs(value)
	if len(dirs) == 0 {
		return nil, nil
	}
	fileConfigs, err := c.loadBinConfigs(dirs[0])
	if len(dirs) == 1 || err != nil {
		return fileConfigs, err
	}

	for i := 1; i < len(dirs); i++ {
		dir := dirs[i]
		cfg, err := c.loadBinConfigs(dir)
		if err != nil {
			return nil, err
		}
		fileConfigs = arrmap.MergeSlices(fileConfigs, cfg)
	}

	return fileConfigs, nil
}
func (c *Config) loadBinConfigs(dir string) (map[string][]byte, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read bin config directory %s: %w", dir, err)
	}

	rsaFiles := make(map[string][]byte, len(entries))
	// 因为file config是配对出现的，所以要整体加载
	for _, entry := range entries {
		if isNotValidFile(entry) {
			continue
		}
		if err = c.loadFile(dir, entry, rsaFiles); err != nil {
			return nil, err
		}
	}
	if len(rsaFiles) > 0 {
		return rsaFiles, nil
	}
	return nil, nil
}

// loadFile 加载单个RSA文件
func (c *Config) loadFile(root string, entry fs.DirEntry, rsaFiles map[string][]byte) error {
	filePath := filepath.Join(root, entry.Name())

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read bin config file %s: %w", filePath, err)
	}

	if len(data) == 0 {
		return fmt.Errorf("bin config file is empty: %s", filePath)
	}

	rsaFiles[entry.Name()] = data
	return nil
}

// 不要获取太细分，否则容易导致错误不容易被排查
func (c *Config) getBinData(name string) []byte {
	if c.isOnWrite() {
		cfgMtx.RLock()
		defer cfgMtx.RUnlock()
	}
	if len(c.binConfig) == 0 {
		return nil
	}
	return c.binConfig[name]
}
func (c *Config) isOnWrite() bool {
	return c.onWrite.Load()
}
