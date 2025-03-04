package alog

import (
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"log"
)

// Console Printc message to docker/console debugging, and save it in the log file
// styles: afmt Colors
func Console(msg string, styles ...string) {
	// for docker or console debugging
	afmt.Console(msg, styles...)
	// log file is no way to display color
	log.Println(msg)
}

func Log(msg string, a ...any) {
	Console(fmt.Sprintf(msg, a...))
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
