package boot

import (
	"context"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/aconfig"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
	"os"
	"os/signal"
	"syscall"
)

var (
	sigs = make(chan os.Signal, 1)
)

func listenTerminateSignals(app *aa.App) {
	// SIGINT: Ctrl + C; SIGTERM: shutdown or container stopped
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		app.Log.Warnf(app.GlobalContext, "terminate signal: %d", sig)
		app.GlobalCancel()
		driver.CloseAllPools(nil)
	}()
}

func InitApp(ctx acontext.Context, cancel context.CancelFunc, configPath string, selfTest bool, profile *debug.Profile) *aa.App {
	profile = profile.Fork("init app openlab").WithLabel("boot")
	cfg, err := aconfig.New(configPath, configValueProcessor)
	ae.PanicOnErrs(err)

	logLevel := alog.NameToLevel(cfg.GetString("log_level", "debug"))
	app := aa.New(ctx, cancel, cfg).WithLog(alog.NewDefaultLog(logLevel))
	redirectLog(app)
	app.Config.Log()
	// loadOtherConfigs(app)

	if selfTest {
		SelfTest(app)
	}

	register(app)
	listenTerminateSignals(app)
	return app
}

func redirectLog(app *aa.App) {
	cfg := app.Config
	dir := cfg.GetString("app.log_dir")
	logBufferSize := cfg.Get("app.log_buffer_size").DefaultInt(0)
	logSymlink := cfg.GetString("app.log_symlink")
	err := debug.RedirectLog(dir, 0666, logBufferSize, logSymlink)
	if err != nil {
		panic(err.Error())
	}
}

func configValueProcessor(key string, value string) (string, error) { return value, nil }