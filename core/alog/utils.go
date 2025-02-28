package alog

import (
	"fmt"
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"log"
)

// Println Print message to docker/console debugging, and save it in the log file
// styles: afmt Colors
func Println(msg string, styles ...string) {
	// for docker or console debugging
	afmt.Println(msg, styles...)
	// log file is no way to display color
	log.Println(msg)
}

func Log(msg string, a ...any) {
	msg = fmt.Sprintf(msg+"\n", a...)
	Println(msg)
}

func PrintError(err error) {
	if err != nil {
		Println(err.Error())
	}
}

func PrintE(e *ae.Error) {
	if e != nil {
		Println(e.String())
	}
}
