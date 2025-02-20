package utils

import (
	"io"
	"log"
)

func Close(closers ...io.Closer) {
	if len(closers) == 0 {
		return
	}
	for _, closer := range closers {
		if closer == nil {
			continue
		}
		if err := closer.Close(); err != nil {
			log.Printf("error closing %v: %v", closer, err)
		}
	}
}
