// +build debug

package main

import "fmt"

func debugf(format string, a ...interface{}) {
	log(fmt.Sprintf(format, a...))
}
