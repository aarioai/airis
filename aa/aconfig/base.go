package aconfig

import (
	"github.com/hashicorp/consul/api"
	"io/fs"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Environment constants
const (
	EnvLocal       Env = "local"
	EnvDevelopment Env = "development"
	EnvProduction  Env = "production"
	EnvTesting     Env = "testing"
	EnvStaging     Env = "staging"
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
	envs   = []string{EnvLocal.String(), EnvDevelopment.String(), EnvProduction.String(), EnvTesting.String(), EnvStaging.String()}
	cfgMtx sync.RWMutex
)

type Env string

// NewEnv
// Example NewEnv("testing") NewEnv("testing_project") NewEnv("project_testing")
func NewEnv(env string) Env {
	// keep it, faster than check has prefix and has suffix
	if slices.Contains(envs, env) {
		return Env(env)
	}
	for _, environment := range envs {
		if strings.HasPrefix(env, environment+"_") || strings.HasSuffix(env, "_"+environment) {
			return Env(environment)
		}
	}
	return Env(env)
}

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
	textConfig  map[string]string // 使用字符串保存，避免[]byte被业务层修改
	otherConfig map[string]string // 不要使用 sync.Map， 直接对整个map加锁设置
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

func (c *Config) SetConsul(name string, client *api.Client) {
	cfgMtx.Lock()
	defer cfgMtx.Unlock()
	c.consulMap[name] = client
}

func (c *Config) SetDefaultConsul(client *api.Client) {
	cfgMtx.Lock()
	defer cfgMtx.Unlock()
	c.consulMap["DEFAULT"] = client
}

func (c *Config) Consul(name string) *api.Client {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()
	return c.consulMap[name] // panic on doesn't exist
}

func (c *Config) DefaultConsul() *api.Client {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()
	return c.consulMap["DEFAULT"] // panic on doesn't exist
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
