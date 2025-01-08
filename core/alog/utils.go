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
func Debug(msg string) {
	afmt.Println("[debug] "+msg, afmt.Cyan)
}
func Info(msg string) {
	afmt.Println("[info] "+msg, afmt.Yellow)
}
func Notice(msg string) {
	afmt.Println("[notice] "+msg, afmt.Magenta)
}
func Warn(msg string) {
	afmt.Println("[warn] "+msg, afmt.Yellow)
}
func Error(msg string) {
	afmt.Println("[error] "+msg, afmt.Red)
}
