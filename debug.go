package main

import (
	"fmt"
)

func debugf(format string, a ...interface{}) {
	if verbose {
		log(fmt.Sprintf("// "+format, a...))
	}
}
