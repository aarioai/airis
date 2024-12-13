package config

import (
	"github.com/aarioai/airis/core/atype"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type MysqlPoolConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// https://github.com/go-sql-driver/mysql/
type MysqlConfig struct {
	Schema   string // dbname
	User     string
	Password string
	// Scheme   string // tcp|unix，只支持tcp，unix仅本地可用
	TLS  string // 默认 false，Valid Values:   true, false, skip-verify, preferred, <name>
	Host string
	// Charset  string  不建议用，应该服务器默认设置

	// mysql客户端在尝试与mysql服务器建立连接时，mysql服务器返回错误握手协议前等待客户端数据包的最大时限。默认10秒。
	ConnTimeout  time.Duration // 使用时，需要设置单位，s, ms等。Timeout for establishing connections, aka dial timeout
	ReadTimeout  time.Duration // 使用时，需要设置单位，s, ms等。I/O read timeout.
	WriteTimeout time.Duration // 使用时，需要设置单位，s, ms等。I/O write timeout.

	Pool MysqlPoolConfig
}

// ParseTimeout connection timeout, r timeout, w timeout, heartbeat interval
// 10s, 1000ms
func (c *Config) ParseTimeout(t string, defaultTimeouts ...time.Duration) (conn time.Duration, read time.Duration, write time.Duration) {
	for i, t := range defaultTimeouts {
		switch i {
		case 0:
			conn = t
		case 1:
			read = t
		case 2:
			write = t
		}
	}

	ts := strings.Split(strings.Replace(t, " ", "", -1), ",")
	for i, t := range ts {
		switch i {
		case 0:
			conn = parseToDuration(t)
		case 1:
			read = parseToDuration(t)
		case 2:
			write = parseToDuration(t)
		}
	}

	return
}
func (c *Config) tryGetMysqlCfg(section string, key string) (string, error) {
	k := section + "." + key
	v, err := c.MustGetString(k)
	if err == nil {
		return v, nil
	}
	if section != "mysql" {
		if !strings.HasPrefix(section, "mysql_") {
			// 尝试section加 mysql_ 开头
			return c.tryGetMysqlCfg("mysql_"+section, key)
		}
		// 读取默认值，即 mysql.$key
		return c.MustGetString("mysql." + key)
	}
	return "", err
}
func (c *Config) MysqlConfig(section string) (MysqlConfig, error) {
	if section == "" {
		section = "mysql"
	}
	host, err := c.tryGetMysqlCfg(section, "host")
	if err != nil {
		return MysqlConfig{}, err
	}
	schema, err := c.tryGetMysqlCfg(section, "schema")
	if err != nil {
		// schema 如果不存在，那么就跟section保持一致
		schema = section
	}
	user, err := c.tryGetMysqlCfg(section, "user")
	if err != nil {
		return MysqlConfig{}, err
	}
	password, err := c.tryGetMysqlCfg(section, "password")
	if err != nil {
		return MysqlConfig{}, err
	}

	tls, _ := c.tryGetMysqlCfg(section, "tls")
	timeout, _ := c.tryGetMysqlCfg(section, "timeout")
	ct, rt, wt := c.ParseTimeout(timeout)
	poolMaxIdleConns, _ := c.tryGetMysqlCfg(section, "pool_max_idle_conns")
	poolMaxOpenConns, _ := c.tryGetMysqlCfg(section, "pool_max_open_conns")
	poolConnMaxLifetime, _ := c.tryGetMysqlCfg(section, "pool_conn_max_life_time")
	poolConnMaxIdleTime, _ := c.tryGetMysqlCfg(section, "pool_conn_max_idle_time")

	newV := atype.New()
	defer newV.Release()

	cf := MysqlConfig{
		Schema:       schema,
		User:         user,
		Password:     password,
		TLS:          tls,
		Host:         host,
		ConnTimeout:  ct,
		ReadTimeout:  rt,
		WriteTimeout: wt,
		Pool: MysqlPoolConfig{
			MaxIdleConns:    newV.Reload(poolMaxIdleConns).DefaultInt(0),
			MaxOpenConns:    newV.Reload(poolMaxOpenConns).DefaultInt(0),
			ConnMaxLifetime: time.Duration(newV.Reload(poolConnMaxLifetime).DefaultInt64(0)) * time.Second,
			ConnMaxIdleTime: time.Duration(newV.Reload(poolConnMaxIdleTime).DefaultInt64(0)) * time.Second,
		},
	}
	return cf, nil
}

func (c *Config) tryGeRedisCfg(section string, key string) (string, error) {
	k := section + "." + key
	v, err := c.MustGetString(k)
	if err == nil {
		return v, nil
	}
	if section != "redis" {
		if !strings.HasPrefix(section, "redis_") {
			// 尝试section加 redis_ 开头
			return c.tryGeRedisCfg("redis_"+section, key)
		}

		// 读取默认值，即 redis.$key
		v, err = c.MustGetString("redis." + key)
		if err == nil {
			return v, nil
		}
	}
	return "", err
}

func (c *Config) RedisConfig(section string) (*redis.Options, error) {
	if section == "" {
		section = "redis"
	}

	addr, err := c.tryGeRedisCfg(section, "addr")
	if err != nil {
		return nil, err
	}
	network, _ := c.tryGeRedisCfg(section, "network")
	clientName, _ := c.tryGeRedisCfg(section, "client_name")
	protocol, _ := c.tryGeRedisCfg(section, "protocol")
	username, _ := c.tryGeRedisCfg(section, "username") // username 可以为空
	password, _ := c.tryGeRedisCfg(section, "password") // password 可以为空
	db, _ := c.tryGeRedisCfg(section, "db")
	maxRetries, _ := c.tryGeRedisCfg(section, "max_retries")
	minRetryBackoff, _ := c.tryGeRedisCfg(section, "min_retry_backoff")
	maxRetryBackoff, _ := c.tryGeRedisCfg(section, "max_retry_backoff")
	dialTimeout, _ := c.tryGeRedisCfg(section, "dial_timeout")
	readTimeout, _ := c.tryGeRedisCfg(section, "read_timeout")
	writeTimeout, _ := c.tryGeRedisCfg(section, "write_timeout")
	contextTimeoutEnabled, _ := c.tryGeRedisCfg(section, "context_timeout_enabled")
	poolFIFO, _ := c.tryGeRedisCfg(section, "pool_fifo")
	poolSize, _ := c.tryGeRedisCfg(section, "pool_size")
	poolTimeout, _ := c.tryGeRedisCfg(section, "pool_timeout")
	minIdleConns, _ := c.tryGeRedisCfg(section, "min_idle_conns")
	maxIdleConns, _ := c.tryGeRedisCfg(section, "max_idle_conns")
	maxActiveConns, _ := c.tryGeRedisCfg(section, "max_active_conns")
	connMaxIdleTime, _ := c.tryGeRedisCfg(section, "conn_max_idle_time")
	connMaxLifetime, _ := c.tryGeRedisCfg(section, "conn_max_lifetime")
	disableIdentity, _ := c.tryGeRedisCfg(section, "disable_identity")
	identitySuffix, _ := c.tryGeRedisCfg(section, "identity_suffix")
	unstableResp3, _ := c.tryGeRedisCfg(section, "unstable_resp3")
	newV := atype.New()
	defer newV.Release()

	opt := redis.Options{
		Network: network,
		Addr:    addr, //  127.0.0.1:6379
		// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
		ClientName:                 clientName,
		Dialer:                     nil,
		OnConnect:                  nil,
		Protocol:                   atype.ParseInt(protocol),
		Username:                   username,
		Password:                   password,
		CredentialsProvider:        nil,
		CredentialsProviderContext: nil,
		DB:                         atype.ParseInt(db),
		MaxRetries:                 atype.ParseInt(maxRetries),
		MinRetryBackoff:            time.Duration(atype.ParseInt64(minRetryBackoff)),
		MaxRetryBackoff:            time.Duration(atype.ParseInt64(maxRetryBackoff)),
		DialTimeout:                time.Duration(atype.ParseInt64(dialTimeout)),
		ReadTimeout:                time.Duration(atype.ParseInt64(readTimeout)),
		WriteTimeout:               time.Duration(atype.ParseInt64(writeTimeout)),
		ContextTimeoutEnabled:      atype.ParseBool(contextTimeoutEnabled),
		PoolFIFO:                   atype.ParseBool(poolFIFO),
		PoolSize:                   atype.ParseInt(poolSize),
		PoolTimeout:                time.Duration(atype.ParseInt64(poolTimeout)),
		MinIdleConns:               atype.ParseInt(minIdleConns),
		MaxIdleConns:               atype.ParseInt(maxIdleConns),
		MaxActiveConns:             atype.ParseInt(maxActiveConns),
		ConnMaxIdleTime:            time.Duration(atype.ParseInt64(connMaxIdleTime)),
		ConnMaxLifetime:            time.Duration(atype.ParseInt64(connMaxLifetime)),
		TLSConfig:                  nil,
		Limiter:                    nil,
		// 官方写错，会在 v10 更正过来
		DisableIndentity: atype.ParseBool(disableIdentity),
		IdentitySuffix:   identitySuffix,
		UnstableResp3:    atype.ParseBool(unstableResp3),
	}
	return &opt, nil
}
