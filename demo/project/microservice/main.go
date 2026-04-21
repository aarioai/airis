package main

import (
	"flag"
	"fmt"
	"project/microservice/boot"
	"runtime"

	"github.com/aarioai/airis/aa/alog"
)

var (
	configPath = flag.String("config", "./config/app.ini", "config path")
	alt        = flag.Bool("alt", false, "switch to alternated function")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in main", r) // for docker container debug
		}
	}()

	alog.Printf("config: %s, alt: %v", *configPath, *alt)
	boot.Boot(*configPath, *alt)
}
