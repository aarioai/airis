package main

import (
	"flag"
	"fmt"
	"project/simple/boot"
	"runtime"

	"github.com/aarioai/airis/aa/alog"
)

var (
	configPath = flag.String("config", "./config/app_simple.ini", "config path")
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
