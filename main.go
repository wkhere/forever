package main // import "github.com/wkhere/forever"

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type configT struct {
	dir     string
	minTick time.Duration
	verbose bool

	progConfig progConfigT
}

func parseArgs() (c *configT) {
	c = new(configT)

	var helpOnly bool
	flagset := pflag.NewFlagSet("flags", pflag.ContinueOnError)
	flagset.SortFlags = false

	flagset.StringVarP(&c.dir, "dir", "d", ".",
		"switch to `directory`")
	flagset.DurationVarP(&c.minTick, "tick", "t", 200*time.Millisecond,
		"events tick")
	flagset.BoolVarP(&c.verbose, "verbose", "v", false,
		"be verbose")
	flagset.BoolVarP(&helpOnly, "help", "h", false,
		"show this help and exit")

	err := flagset.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "forever:", err)
		os.Exit(2)
	}
	if helpOnly {
		help(flagset)
		os.Exit(0)
	}

	setupDebug(c.verbose)

	if rest := flagset.Args(); len(rest) > 0 {
		c.progConfig.explicitProg = true
		c.progConfig.prog = rest[0]
		c.progConfig.args = rest[1:]
	}
	return
}

func help(f *pflag.FlagSet) {
	f.SetOutput(os.Stdout)
	p := fmt.Printf
	p("Usage: forever [-d dir] [-t events-tick] [-v] [program...]\n\n")
	f.PrintDefaults()
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
	select {} // block forever
}
