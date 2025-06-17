package aconfig

import (
	"slices"
	"strings"
)

// Environment constants
const (
	EnvLocal       Env = "local"
	EnvDevelopment Env = "development"
	EnvTest        Env = "test"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

var (
	envs        = []string{EnvLocal.String(), EnvDevelopment.String(), EnvTest.String(), EnvStaging.String(), EnvProduction.String()}
	envAliasMap = map[Env][]string{
		EnvDevelopment: {"develop", "dev"},
		EnvTest:        {"testing"},
		EnvStaging:     {"stag"},
		EnvProduction:  {"product", "prod", "pro"},
	}
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
	for v, aliases := range envAliasMap {
		for _, alias := range aliases {
			if env == alias {
				return v
			}
			if strings.HasPrefix(env, alias+"_") {
				return NewEnv(EnvDevelopment.String() + strings.TrimPrefix(env, alias))
			}
		}
	}

	return Env(env)
}

func (env Env) String() string          { return string(env) }
func (env Env) IsLocal() bool           { return env == EnvLocal }
func (env Env) IsDevelopment() bool     { return env == EnvDevelopment }
func (env Env) IsTest() bool            { return env == EnvTest }
func (env Env) IsStaging() bool         { return env == EnvStaging }
func (env Env) IsProduction() bool      { return env == EnvProduction }
func (env Env) BeforeDevelopment() bool { return env.IsLocal() }
func (env Env) BeforeTest() bool        { return env.IsDevelopment() || env.BeforeDevelopment() }
func (env Env) BeforeStaging() bool     { return env.IsTest() || env.BeforeTest() }
func (env Env) BeforeProduction() bool  { return env.IsStaging() || env.BeforeStaging() }
func (env Env) AfterStaging() bool      { return env.IsProduction() }
func (env Env) AfterTest() bool         { return env.IsStaging() || env.AfterStaging() }
func (env Env) AfterDevelopment() bool  { return env.IsTest() || env.AfterTest() }
func (env Env) AfterLocal() bool        { return env.IsDevelopment() || env.AfterDevelopment() }
