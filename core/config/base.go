package config

import (
	"fmt"
	"github.com/aarioai/airis/pkg/utils"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// Environment constants
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvTesting     = "testing"
	EnvStaging     = "staging"
)

// Configuration keys
const (
	CkRsaRoot    = "rsa_root"
	CkEnv        = "env"
	CkTimezoneID = "timezone_id"
	CkTimeFormat = "time_format"
	CkMock       = "mock"
)

type Env string

func (env Env) String() string      { return string(env) }
func (env Env) IsDevelopment() bool { return env == EnvDevelopment }
func (env Env) IsTesting() bool     { return env == EnvTesting }
func (env Env) IsStaging() bool     { return env == EnvStaging }
func (env Env) IsProduction() bool  { return env == EnvProduction }
func (env Env) BeforeTesting() bool { return env.IsTesting() || env.IsDevelopment() }
func (env Env) BeforeStaging() bool { return env.IsStaging() || env.BeforeTesting() }
func (env Env) AfterStaging() bool  { return env.IsStaging() || env.IsProduction() }
func (env Env) AfterTesting() bool  { return env.IsTesting() || env.AfterStaging() }

type Snapshot struct {
	data        map[string]string
	rsa         map[string][]byte // rsa 是配对出现的，不要使用 sync.Map， 直接对整个map加锁设置
	otherConfig map[string]string // 不要使用 sync.Map， 直接对整个map加锁设置
}

type Config struct {
	/*
		https://en.wikipedia.org/wiki/Deployment_environment
		local -> development/trunk -> integration -> testing/test/qc/internal acceptnace -> staging/stage/model/pre-production/demo -> production/live
		development -> test -> pre-production -> production
	*/
	Env          Env
	TimezoneID   string // e.g. "Asia/Shanghai"
	TimeLocation *time.Location
	TimeFormat   string // e.g. "2006-02-01 15:04:05"
	Mock         bool   // using mock

	onWrite     atomic.Bool
	path        string
	data        map[string]string
	rsa         map[string][]byte // rsa 是配对出现的，不要使用 sync.Map， 直接对整个map加锁设置
	otherConfig map[string]string // 不要使用 sync.Map， 直接对整个map加锁设置
	snapshot    atomic.Pointer[Snapshot]
}

func (c *Config) Log() {
	info := fmt.Sprintf(`
Launch Configuration:
Environment: %s
Timezone: %s
Mock Enabled: %v
Git Version: %s
`,
		c.Env,
		c.TimezoneID,
		c.Mock,
		utils.GitVersion(),
	)

	// 方便运行程序时直接显示
	fmt.Println(info)
	// 记录进日志，方便通过消息队列通知
	log.Println(info)
}

func parseToDuration(d string) time.Duration {
	if len(d) < 2 {
		return 0
	}
	var t int
	if d[len(d)-2:] == "ms" {
		t, _ = strconv.Atoi(d[0 : len(d)-2])
		return time.Duration(t) * time.Millisecond
	}

	if d[len(d)-1:] == "s" {
		t, _ = strconv.Atoi(d[0 : len(d)-1])
	} else {
		t, _ = strconv.Atoi(d)
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
