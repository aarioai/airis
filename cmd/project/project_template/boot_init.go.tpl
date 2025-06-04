package boot

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/aconfig"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
)

func InitApp(ctx acontext.Context, configPath string, selfTest bool, profile *debug.Profile) *aa.App {
	profile = profile.Fork("init app {{APP_NAME}}").WithLabel("boot")
	cfg, err := aconfig.New(configPath, configValueProcessor)
	ae.PanicOnErrs(err)

	logLevel := alog.NameToLevel(cfg.GetString("log_level", "debug"))
	app := aa.New(cfg, alog.NewDefaultLog(logLevel))
	app.Config.Log()
	// loadOtherConfigs(app)
	redirectLog(app)

	if selfTest {
		SelfTest(app)
	}

	register(app)
	return app
}

func configValueProcessor(key string, value string) (string, error) {
	return value, nil
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
