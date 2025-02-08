package aconfig_test

import (
	"github.com/aarioai/airis/core/aconfig"
	"testing"
	"time"
)

func TestParseIni(t *testing.T) {
	c, err := aconfig.New("./parse_ini_test.ini", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Extend(map[string]string{
		"default.time": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		t.Fatal(err)
	}

	//c.Dump()
	//c.Log()

	debug, err := c.Get("debug").Bool()
	if err != nil || !debug {
		t.Errorf("config parse debug fail: %s", c.Get("debug").String())
	}
	
}
