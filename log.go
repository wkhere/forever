package main

import (
	"fmt"
	"os"
)

var logw = os.Stderr

func log(a ...any) { fmt.Fprintln(logw, a...) }

func logf(format string, a ...any) {
	fmt.Fprintf(logw, format, a...)
	fmt.Fprintln(logw)
}

func logfColor(color int, format string, a ...any) {
	fmt.Fprintf(logw, "\033[3%d;1m", color)
	fmt.Fprintf(logw, format, a...)
	fmt.Fprint(logw, "\033[0m\n")
}

func logfBlue(format string, a ...any) {
	logfColor(4, format, a...)
}

func logfGreen(format string, a ...any) {
	logfColor(2, format, a...)
}

func logfRed(format string, a ...any) {
	logfColor(1, format, a...)
}

var (
	debug  = func(a ...any) {}
	debugf = func(format string, a ...any) {}

	watchdebug = func(format string, a ...any) {}
)

func setupDebug(ok bool) {
	if ok {
		debug = func(a ...any) {
			fmt.Fprint(logw, "// ")
			fmt.Fprintln(logw, a...)
		}
		debugf = func(format string, a ...any) {
			fmt.Fprint(logw, "// ")
			fmt.Fprintf(logw, format, a...)
			fmt.Fprintln(logw)
		}
		watchdebug = _watchdebug
	}
}
