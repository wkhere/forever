package main // import "github.com/wkhere/forever"

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type configT struct {
	dir     string
	minTick time.Duration
	verbose bool
	debug   bool

	progConfig progConfigT
}

func parseArgs() (c *configT) {
	c = new(configT)

	flag.StringVar(&c.dir, "d", ".", "directory")
	flag.DurationVar(&c.minTick, "t", 200*time.Millisecond, "events tick")
	flag.BoolVar(&c.verbose, "v", false, "verbose mode")
	flag.BoolVar(&c.debug, "vv", false, "debug mode")

	flag.Usage = usage
	flag.Parse()

	setupDebug(c.debug)

	if rest := flag.Args(); len(rest) > 0 {
		c.progConfig.explicitProg = true
		c.progConfig.prog = rest[0]
		c.progConfig.args = rest[1:]
	}
	return
}

func usage() {
	p := func(format string, a ...interface{}) {
		fmt.Fprintf(flag.CommandLine.Output(), format, a...)
	}
	p("Usage: forever [-d dir] [-t events-tick] [-v|-vv] [program...]\n")
	flag.PrintDefaults()
	p("\nIf program is not given, the following will be tried:\n\t%s\n",
		strings.Join(defaultProgs, "\n\t"),
	)
	if writeDirsOnSignal {
		p("\nThe list of watched directories can be dumped to a file")
		p("\n%s by sending HUP (-1) signal.\n", writeDirsOutputPattern)
	}
}

func main() {
	config := parseArgs()

	err := os.Chdir(config.dir)
	if err != nil {
		fatal("cannot prepare:", err)
	}

	w, err := newWatcher(config.minTick)
	if err != nil {
		fatal("cannot start watcher:", err)
	}

	// watcher should add all files before looping
	w.feed()
	w.installSignal()

	go loop(w, &config.progConfig)

	neverending := make(chan struct{})
	<-neverending
}
