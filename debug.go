package main

import (
	"flag"
	"fmt"
)

var debugFlag = flag.Bool("d", false, "debug mode")

func debugf(format string, a ...interface{}) {
	if *debugFlag {
		log(fmt.Sprintf(format, a...))
	}
}
