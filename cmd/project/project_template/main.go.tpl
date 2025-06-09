package main

import (
	"flag"
	"fmt"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
	"{{ROOT}}/boot"
	"{{ROOT}}/router"
	"runtime"
)

var (
	configPath  = flag.String("config", "./config/app-local.ini", "config path")
	selfTest    = flag.Bool("selftest", false, "self test")
	ctx, cancel = acontext.WithCancel(acontext.Background())
	profile     = debug.DefaultProfile()
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in main", r) // for docker container log
		}
	}()

	alog.Printf("config: %s, self test: %v", *configPath, *selfTest)
	app := boot.InitApp(ctx, cancel, *configPath, *selfTest, profile)
	router.Run(app, profile)
}
