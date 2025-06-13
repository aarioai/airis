package main

import (
	"flag"
	"fmt"
	"github.com/aarioai/airis/aa/alog"
	"{{ROOT}}/boot"
	"runtime"
)

var (
	configPath  = flag.String("config", "./config/app-local.ini", "config path")
	selfTest    = flag.Bool("selftest", false, "self test")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in main", r) // for docker container debug
		}
	}()

	alog.Printf("config: %s, self test: %v", *configPath, *selfTest)

	boot.Boot(*configPath, *selfTest)
}
