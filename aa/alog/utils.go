package alog

import (
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"log"
)

var (
	StartPrefix = "start: "
	StopPrefix  = "stop: "
)

// Console Printc message to docker/console debugging, and save it in the log file
// styles: afmt Colors
func Console(msg string, styles ...string) {
	// for docker or console debugging
	afmt.Console(msg, styles...)
	// log file is no way to display color
	log.Println(msg)
}

func Log(msg string, args ...any) {
	Console(fmt.Sprintf(msg, args...))
}

func LogOnError(err error) {
	if err != nil {
		Console(err.Error())
	}
}

func LogOn(e *ae.Error) {
	if e != nil {
		Console(e.String())
	}
}

func Alerting(msg string, args ...any) {
	Console(afmt.Sprintf(msg, args...), afmt.Red)
}

func Warning(msg string, args ...any) {
	Console(afmt.Sprintf(msg, args...), afmt.Yellow)
}

func Information(msg string, args ...any) {
	Console(afmt.Sprintf(msg, args...), afmt.Magenta)
}

func Start(name string, args ...any) {
	name = afmt.Sprintf(name, args...)
	Console(StartPrefix+name, afmt.Green)
}

func Stop(name string, args ...any) {
	name = afmt.Sprintf(name, args...)
	Console(StopPrefix+name, afmt.Red)
}
