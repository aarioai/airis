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
	var (
		name string
		err  error
		dir  []fs.DirEntry
		p    string // file path
	)

	if dir, err = os.ReadDir(root); err != nil {
		return err
	}

	rsas := make(map[string][]byte, len(dir))
	// 因为RSA是配对出现的，所以要整体加载
	for _, entry := range dir {
		name = entry.Name()
		// 跳过目录和隐藏文件
		if entry.IsDir() || len(name) == 0 || name[0] == '.' {
			continue
		}

		p = filepath.Join(root, name)
		dat, err := os.ReadFile(p)
		if err != nil || len(dat) == 0 {
			return fmt.Errorf("invalid rsa file `%s`: %w", p, err)
		}
		rsas[name] = dat

	}
	if len(rsas) > 0 {
		c.AddRsaConfigs(rsas)
	}
	return nil
}

func (c *Config) AddRsaConfigs(rsaConfigs map[string][]byte) {
	cfgMtx.Lock()
	cfgMtx.Unlock()
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
	v, ok := c.rsa[name]
	return v, ok
}
