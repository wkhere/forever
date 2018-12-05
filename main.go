package main // import "github.com/wkhere/forever"

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type config struct {
	dir     string
	minTick time.Duration
	verbose bool
}

// vars which need to be global
var (
	verbose bool
)

func parseArgs() (c config) {

	flag.DurationVar(&c.minTick, "t", 200*time.Millisecond, "events tick")
	flag.BoolVar(&c.verbose, "v", false, "verbose/debug mode")

	flag.Usage = usage
	flag.Parse()

	c.dir = "." //tmp

	return
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: forever [-v] [-t events-tick]\n")
	flag.PrintDefaults()
}

func main() {
	c := parseArgs()
	verbose = c.verbose

	err := os.Chdir(c.dir)
	if err != nil {
		log("cannot prepare:", err)
		os.Exit(1)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log("cannot start watcher: ", err)
		os.Exit(1)
	}

	// watcher should add all files before looping
	feedWatcher(w, c.dir)

	go loop(w, c.minTick)

	neverending := make(chan struct{})
	<-neverending
}

func log(msgs ...interface{}) {
	fmt.Fprintln(os.Stderr, msgs...)
}

func logf(format string, msgs ...interface{}) {
	log(fmt.Sprintf(format, msgs...))
}

func logBlue(s string) {
	io.WriteString(os.Stderr, "\033[34m")
	io.WriteString(os.Stderr, s)
	io.WriteString(os.Stderr, "\033[0m\n")
}
