package config

import (
	"fmt"
	"github.com/aarioai/airis/pkg/utils"
	"gopkg.in/ini.v1"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvTesting     = "testing"
	EnvStaging     = "staging"

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

	path        string
	data        *ini.File
	rsa         map[string][]byte // rsa 是配对出现的，不要使用 sync.Map， 直接对整个map加锁设置
	otherConfig map[string]string // 不要使用 sync.Map， 直接对整个map加锁设置
}

func (c *Config) Log() {
	msg := fmt.Sprintf("lauching...\nenv: %s\ntimezone_id: %s\nmock: %v\ngit_ver: %s", c.Env, c.TimezoneID, c.Mock, utils.GitVersion())
	// print on console
	fmt.Println(msg)
	log.Println(msg)
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

func splitDots(keys ...string) []string {
	n := make([]string, 0)
	for _, key := range keys {
		n = append(n, strings.Split(key, ".")...)
	}
	return n
}
