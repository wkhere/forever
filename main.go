package main // import "github.com/wkhere/forever"

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type configT struct {
	dir     string
	minTick time.Duration
	verbose bool
	debug   bool

	progConfig progConfigT
}

var config *configT

func parseArgs() (c *configT) {
	c = new(configT)

	flag.StringVar(&c.dir, "d", ".", "directory")
	flag.DurationVar(&c.minTick, "t", 200*time.Millisecond, "events tick")
	flag.BoolVar(&c.verbose, "v", false, "verbose mode")
	flag.BoolVar(&c.debug, "vv", false, "debug mode")

	flag.Usage = usage
	flag.Parse()

	if rest := flag.Args(); len(rest) > 0 {
		c.progConfig.explicitProg = true
		c.progConfig.prog = rest[0]
		c.progConfig.args = rest[1:]
	}
	return
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: forever [-d dir] [-t events-tick] [-v|-vv] [command]\n")
	flag.PrintDefaults()
}

func main() {
	config = parseArgs()

	err := os.Chdir(config.dir)
	if err != nil {
		fatal("cannot prepare:", err)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		fatal("cannot start watcher: ", err)
	}

	// watcher should add all files before looping
	feedWatcher(w)

	go loop(w, config.minTick, &config.progConfig)

	neverending := make(chan struct{})
	<-neverending
}
