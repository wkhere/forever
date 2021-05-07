package main

import (
	logpkg "log"

	"fmt"
	"io/ioutil"
	"os"
)

var (
	logger      = logpkg.New(os.Stderr, "", 0)
	debugLogger = logpkg.New(os.Stderr, "// ", 0)

	log   = logger.Println
	logf  = logger.Printf
	fatal = logger.Fatalln

	debugf = debugLogger.Printf
)

func logfBlue(format string, a ...interface{}) {
	logfColor(4, format, a...)
}

func logfGreen(format string, a ...interface{}) {
	logfColor(2, format, a...)
}

func logfRed(format string, a ...interface{}) {
	logfColor(1, format, a...)
}

func logfColor(color int, format string, a ...interface{}) {
	w := logger.Writer()
	fmt.Fprintf(w, "\033[3%d;1m", color)
	fmt.Fprintf(w, format, a...)
	fmt.Fprintf(w, "\033[0m\n")
}

func setupDebug(ok bool) {
	if !ok {
		debugLogger.SetOutput(ioutil.Discard)
	}
}
