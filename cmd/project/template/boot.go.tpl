package boot

import (
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/helpers/debug"
	"os"
	"os/signal"
	"{{PROJECT_BASE}}/router"
	"syscall"
)

var (
	profile     = debug.DefaultProfile()
	sigs        = make(chan os.Signal, 1)
)

func Boot(configPath string, selfTest bool) {
	app := initApp(configPath, selfTest)

	router.Serve(app, profile)

	waitTerminate(app)
}

func waitTerminate(app *aa.App) {
	// SIGINT: Ctrl + C; SIGTERM: shutdown or container stopped
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	app.Log.Warnf(app.GlobalContext, "terminate signal: %d", sig)
	app.GlobalCancel()
	driver.CloseAllPools(nil)
}
