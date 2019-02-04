package main

import (
	"fmt"
)

func debugf(format string, a ...interface{}) {
	if config.debug {
		log("//", fmt.Sprintf(format, a...))
	}
}
