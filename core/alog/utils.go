package alog

import (
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
