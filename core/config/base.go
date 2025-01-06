package config

import (
	"github.com/aarioai/airis/pkg/types"
	"io/fs"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Environment constants
const (
	EnvLocal       = "local"
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvTesting     = "testing"
	EnvStaging     = "staging"
)

// Reserved configuration keys
const (
	CkEnv            = "env"
	CkMock           = "mock"
	CkTextConfigDirs = "text_config_dirs" // 使用空格隔开
	CkTimeFormat     = "time_format"
	CkTimezoneID     = "timezone_id"
)

var (
	cfgMtx sync.RWMutex
)

type Env string

func (env Env) String() string          { return string(env) }
func (env Env) IsLocal() bool           { return env == EnvLocal }
func (env Env) IsDevelopment() bool     { return env == EnvDevelopment }
func (env Env) IsTesting() bool         { return env == EnvTesting }
func (env Env) IsStaging() bool         { return env == EnvStaging }
func (env Env) IsProduction() bool      { return env == EnvProduction }
func (env Env) BeforeDevelopment() bool { return env.IsLocal() || env.IsDevelopment() }
func (env Env) BeforeTesting() bool {
	return env.IsTesting() || env.BeforeDevelopment()
}
func (env Env) BeforeStaging() bool { return env.IsStaging() || env.BeforeTesting() }
func (env Env) AfterStaging() bool  { return env.IsStaging() || env.IsProduction() }
func (env Env) AfterTesting() bool  { return env.IsTesting() || env.AfterStaging() }

type Snapshot struct {
	baseConfig  map[string]string
	textConfig  map[string]string
	otherConfig map[string]string
}

type Config struct {
	/*
		https://en.wikipedia.org/wiki/Deployment_environment
		local -> development/trunk -> integration -> testing/test/qc/internal acceptnace -> staging/stage/model/pre-production/demo -> production/live
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
	textConfig  map[string]string // 使用字符串保存，避免[]byte被业务层修改
	otherConfig map[string]string // 不要使用 sync.Map， 直接对整个map加锁设置
	snapshot    atomic.Pointer[Snapshot]

	valueProcessor func(key string, value string) (string, error)
}

func New(path string, valueProcessor func(key string, value string) (string, error)) (*Config, error) {
	cfg := &Config{
		TimeFormat:     "2006-02-01 15:04:05",
		TimezoneID:     time.Local.String(),
		TimeLocation:   time.Local,
		path:           path,
		valueProcessor: valueProcessor,
	}
	err := cfg.Reload()
	return cfg, err
}

func parseToDuration(d string) time.Duration {
	if len(d) < 2 {
		return 0
	}
	var t int
	if d[len(d)-2:] == "ms" {
		t, _ = types.Atoi(d[0 : len(d)-2])
		return time.Duration(t) * time.Millisecond
	}

	if d[len(d)-1:] == "s" {
		t, _ = types.Atoi(d[0 : len(d)-1])
	} else {
		t, _ = types.Atoi(d)
	}
	return time.Duration(t) * time.Second
}

// splitDots splits dot-separated strings into parts
// @example ["a.b.c", "d.e"] -> ["a", "b", "c", "d", "e"]
func splitDots(keys ...string) []string {
	n := make([]string, 0)
	for _, key := range keys {
		n = append(n, strings.Split(key, ".")...)
	}
	return n
}

// shouldSkipFile 判断是否应该跳过该文件
// 跳过目录和隐藏文件
func isNotValidFile(entry fs.DirEntry) bool {
	name := entry.Name()
	return entry.IsDir() || len(name) == 0 || name[0] == '.'
}

// 转化为数组，多个目录
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
