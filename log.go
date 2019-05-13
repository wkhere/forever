package main

import (
	logpkg "log"

	"io"
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

func logBlue(s string) {
	w := logger.Writer()
	io.WriteString(w, "\033[34m")
	io.WriteString(w, s)
	io.WriteString(w, "\033[0m\n")
}

func setupDebug(ok bool) {
	if !ok {
		debugLogger.SetOutput(ioutil.Discard)
	}
}
