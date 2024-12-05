package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func (c *Config) loadRsa() error {
	root := c.GetString(CkRsaRoot)
	if root == "" {
		return nil
	}

	dir, err := os.ReadDir(root)
	if err != nil {
		return fmt.Errorf("failed to read RSA directory %s: %w", root, err)
	}

	rsaFiles := make(map[string][]byte, len(dir))
	// 因为RSA是配对出现的，所以要整体加载
	for _, entry := range dir {

		if isNotValidFile(entry) {
			continue
		}
		if err = c.loadRsaFile(root, entry, rsaFiles); err != nil {
			return err
		}
	}
	if len(rsaFiles) > 0 {
		c.AddRsaConfigs(rsaFiles)
	}
	return nil
}

// shouldSkipFile 判断是否应该跳过该文件
// 跳过目录和隐藏文件
func isNotValidFile(entry fs.DirEntry) bool {
	name := entry.Name()
	return entry.IsDir() || len(name) == 0 || name[0] == '.'
}

// loadRsaFile 加载单个RSA文件
func (c *Config) loadRsaFile(root string, entry fs.DirEntry, rsaFiles map[string][]byte) error {
	filePath := filepath.Join(root, entry.Name())

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read RSA file %s: %w", filePath, err)
	}

	if len(data) == 0 {
		return fmt.Errorf("RSA file is empty: %s", filePath)
	}

	rsaFiles[entry.Name()] = data
	return nil
}
func (c *Config) AddRsaConfigs(rsaConfigs map[string][]byte) {
	cfgMtx.Lock()
	defer cfgMtx.Unlock()
	if c.rsa == nil {
		c.rsa = make(map[string][]byte)
	}
	for name, v := range rsaConfigs {
		if len(v) > 0 {
			c.rsa[name] = v
		}
	}
}

// 不要获取太细分，否则容易导致错误不容易被排查
func (c *Config) getRsa(name string) ([]byte, bool) {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()
	v, exists := c.rsa[name]
	return v, exists
}
