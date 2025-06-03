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

func OnError(err error) {
	if err != nil {
		Console(err.Error())
	}
}

func LogOnE(e *ae.Error) {
	if e != nil {
		Console(e.String())
	}
}
func Start(name string) {
	Console(StartPrefix+name, afmt.Green)
}
func Startf(format string, args ...any) {
	Start(fmt.Sprintf(format, args...))
}

func Stop(name string) {
	Console(StopPrefix+name, afmt.Red)
}
func Stopf(format string, args ...any) {
	Stop(fmt.Sprintf(format, args...))
}
