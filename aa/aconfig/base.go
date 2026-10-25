package aconfig

import (
	"github.com/hashicorp/consul/api"
	"io/fs"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Reserved configuration keys
const (
	CkEnv            = "env"
	CkMock           = "mock"
	CkTextConfigDirs = "text_config_dirs" // space-separated
	CkTimeFormat     = "time_format"
	CkTimezoneID     = "timezone_id"
)

var (
	cfgMtx sync.RWMutex
)

type Snapshot struct {
	baseConfig  map[string]string
	textConfig  map[string]string
	otherConfig map[string]string
}

type Config struct {
	/*
		https://en.wikipedia.org/wiki/Deployment_environment
		local -> development/trunk -> integration -> testing/test/qc/internal acceptance -> staging/stage/model/pre-production/demo -> production/live
		development -> test -> pre-production -> production
	*/
	Env            Env
	Mock           bool // using mock
	TextConfigDirs []string
	TimeFormat     string // e.g. "2006-02-01 15:04:05"
	TimezoneID     string // e.g. "Asia/Shanghai"
	TimeLocation   *time.Location

	onWrite     atomic.Bool
	path        string
	baseConfig  map[string]string
	textConfig  map[string]string // stored as string so the business layer cannot mutate a []byte
	otherConfig map[string]string // do not use sync.Map; lock the whole map and assign it directly
	consulMap   map[string]*api.Client
	snapshot    atomic.Pointer[Snapshot]

	valueProcessor func(key string, value string) (string, error)
}

func New(path string, valueProcessor func(key string, value string) (string, error)) (*Config, error) {
	cfg := &Config{
		TimeFormat:     "2006-02-01 15:04:05",
		TimezoneID:     time.Local.String(),
		TimeLocation:   time.Local,
		path:           path,
		baseConfig:     make(map[string]string),
		textConfig:     make(map[string]string),
		otherConfig:    make(map[string]string),
		consulMap:      make(map[string]*api.Client),
		valueProcessor: valueProcessor,
	}
	err := cfg.Reload()
	return cfg, err
}

// shouldSkipFile reports whether the file should be skipped.
// Skip directories and hidden files.
func isNotValidFile(entry fs.DirEntry) bool {
	name := entry.Name()
	return entry.IsDir() || len(name) == 0 || name[0] == '.'
}

// Convert a space-separated value into a slice of directories.
func textConfigDirs(value string) []string {
	if value == "" {
		return nil
	}
	dirs := strings.Fields(value)
	if len(dirs) == 0 {
		return nil
	}
	return dirs
}
