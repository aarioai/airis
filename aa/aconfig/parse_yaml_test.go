package aconfig_test

import (
	"testing"
	"time"

	"github.com/aarioai/airis/aa/aconfig"
)

func TestParseYAML(t *testing.T) {
	c, err := aconfig.New("./parse_yaml_test.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Extend(map[string]string{
		"default.time": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		t.Fatal(err)
	}

	debug, err := c.Get("debug").Bool()
	if err != nil || !debug {
		t.Errorf("config parse debug fail: %s", c.Get("debug").String())
	}

	if got := c.GetString("driver.redis.DEFAULT.addr"); got != "luexu.com" {
		t.Errorf("driver.redis.DEFAULT.addr = %q, want luexu.com", got)
	}

	// test overlays addr/db; password inherited from DEFAULT
	if got := c.GetString("driver.redis.test.addr"); got != "https://luexu.com" {
		t.Errorf("driver.redis.test.addr = %q, want https://luexu.com", got)
	}
	if got := c.GetString("driver.redis.test.password"); got != "Luexu.com" {
		t.Errorf("driver.redis.test.password = %q, want Luexu.com (from DEFAULT)", got)
	}
	if got := c.GetString("driver.redis.test.db"); got != "1" {
		t.Errorf("driver.redis.test.db = %q, want 1", got)
	}

	// test2 inherits addr; overlays password/db
	if got := c.GetString("driver.redis.test2.addr"); got != "luexu.com" {
		t.Errorf("driver.redis.test2.addr = %q, want luexu.com (from DEFAULT)", got)
	}
	if got := c.GetString("driver.redis.test2.password"); got != "Aario" {
		t.Errorf("driver.redis.test2.password = %q, want Aario", got)
	}

	// empty map inherits all DEFAULT fields
	if got := c.GetString("driver.mysql.test.host"); got != "luexu.com" {
		t.Errorf("driver.mysql.test.host = %q, want luexu.com (from DEFAULT)", got)
	}
	if got := c.GetString("driver.mysql.DEFAULT.pool_max_idle_conns"); got != "0" {
		t.Errorf("driver.mysql.DEFAULT.pool_max_idle_conns = %q, want 0", got)
	}

	// hello inherits host, drops password via __DELETED__, adds schema
	if got := c.GetString("driver.mysql.hello.host"); got != "luexu.com" {
		t.Errorf("driver.mysql.hello.host = %q, want luexu.com (from DEFAULT)", got)
	}
	if got := c.GetString("driver.mysql.hello.schema"); got != "helloworld" {
		t.Errorf("driver.mysql.hello.schema = %q, want helloworld", got)
	}
	if got := c.GetString("driver.mysql.hello.password"); got != "" {
		t.Errorf("driver.mysql.hello.password = %q, want empty (__DELETED__)", got)
	}
}
