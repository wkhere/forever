package main

import (
	"fmt"
)

func debugf(format string, a ...interface{}) {
	if config.verbose {
		log(fmt.Sprintf("// "+format, a...))
	}
}
