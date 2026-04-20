package main

import (
	"flag"
	"fmt"
	"github.com/aarioai/airis/aa/alog"
	"{{PROJECT_BASE}}/boot"
	"runtime"
)

var (
	configPath  = flag.String("config", "./config/app_{{APP_NAME}}.ini", "config path")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in main", r) // for docker container debug
		}
	}()

	alog.Printf("config: %s", *configPath)

	boot.Boot(*configPath)
}
