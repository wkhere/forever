package main // import "github.com/wkhere/forever"

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

var tickFlag = flag.Uint("t", 200, "events tick [ms]")

func init() {
	flag.Usage = usage
	flag.Parse()

	minTick = time.Duration(*tickFlag * uint(time.Millisecond))
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: forever [-d] [-t events-tick] [dir]\n")
	flag.PrintDefaults()
}

func main() {
	var dir string

	switch a := flag.Args(); len(a) {
	case 0:
		dir = "."
	case 1:
		dir = a[0]
	default:
		log("give only 1 directory as an argument")
		os.Exit(2)
	}

	err := os.Chdir(dir)
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
	feedWatcher(w, ".")

	go loop(w)

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
