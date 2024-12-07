package config_test

import (
	"github.com/aarioai/airis/core/config"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestParseIni(t *testing.T) {
	c := config.New("./parse_ini_test.ini")
	debug, err := c.Get("debug").Bool()
	if err != nil || !debug {
		t.Errorf("config parse debug fail: %s", c.Get("debug").String())
	}
	testRedisOpt := redis.Options{
		Password: "Luexu.com",
		DB:       1,
		Addr:     "https://luexu.com",
	}
	test2RedisOpt := redis.Options{
		Password: "Aario",
		DB:       2,
		Addr:     "luexu.com",
	}
	mysqlTestOpt := config.MysqlConfig{
		Schema:       "test",
		Host:         "luexu.com",
		User:         "Aario",
		Password:     "Luexu.com",
		WriteTimeout: time.Second * 5,
	}
	mysqlHelloOpt := config.MysqlConfig{
		Schema:       "helloworld",
		Host:         "luexu.com",
		User:         "Aario",
		Password:     "Luexu.com",
		WriteTimeout: time.Second * 5,
	}
	testRedisConfig(t, c, "redis_test", testRedisOpt)
	testRedisConfig(t, c, "test", testRedisOpt)
	testRedisConfig(t, c, "redis_test2", test2RedisOpt)
	testMySQLConfig(t, c, "mysql_test", mysqlTestOpt)
	testMySQLConfig(t, c, "test", mysqlTestOpt)
	testMySQLConfig(t, c, "hello", mysqlHelloOpt)
}

func testRedisConfig(t *testing.T, c *config.Config, section string, want redis.Options) {
	test, err := c.RedisConfig(section)
	if err != nil {
		t.Fatal(err.Error())
	}
	if test.Password != want.Password {
		t.Errorf("test redis password %s not match %s", test.Password, want.Password)
	}
	if test.DB != want.DB {
		t.Errorf("test redis db %d not match %d", test.DB, want.DB)
	}
	if test.Addr != want.Addr {
		t.Errorf("test redis addr %s not match %s", test.Addr, want.Addr)
	}

}
func testMySQLConfig(t *testing.T, c *config.Config, want string, suppose config.MysqlConfig) {
	mysqlConfig, err := c.MysqlConfig(want)
	if err != nil {
		t.Fatal(err.Error())
	}
	if mysqlConfig.Host != suppose.Host {
		t.Errorf("test mysql host %s not match %s", mysqlConfig.Host, suppose.Host)
	}
	if mysqlConfig.User != suppose.User {
		t.Errorf("test mysql user %s not match %s", mysqlConfig.User, suppose.User)
	}
	if mysqlConfig.Password != suppose.Password {
		t.Errorf("test mysql password %s not match %s", mysqlConfig.Password, suppose.Password)
	}
	if mysqlConfig.WriteTimeout != suppose.WriteTimeout {
		t.Errorf("test mysql password %d not match %d", mysqlConfig.WriteTimeout, suppose.WriteTimeout)
	}
}
