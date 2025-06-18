package boot

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/aconfig"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
)

var (
	ctx, cancel = acontext.WithCancel(acontext.Background())
)

func initApp(configPath string, selfTest bool) *aa.App {
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