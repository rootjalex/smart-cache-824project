package cache

import (
	"log"
)

const DEBUG = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if DEBUG > 0 {
		log.Printf(format, a...)
	}
	return
}
