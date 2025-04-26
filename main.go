package main // import "github.com/wkhere/forever"

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type config struct {
	dir      string
	timeslot time.Duration
	verbose  bool

	prog prog
}

func parseArgs() (c *config) {
	c = new(config)

	var helpOnly bool
	flagset := pflag.NewFlagSet("flags", pflag.ContinueOnError)
	flagset.SortFlags = false

	flagset.StringVarP(&c.dir, "dir", "d", ".",
		"switch to `directory`")
	flagset.DurationVarP(&c.timeslot, "timeslot", "t", 200*time.Millisecond,
		"timeslot for write events")
	flagset.BoolVarP(&c.verbose, "verbose", "v", false,
		"be verbose")
	flagset.BoolVarP(&helpOnly, "help", "h", false,
		"show this help and exit")

	flagset.Usage = func() {
		p := func(format string, a ...interface{}) {
			fmt.Fprintf(flagset.Output(), format, a...)
		}
		p("Usage: forever [-d dir] [-t duration] [-v] [-- program ...]\n\n")

		flagset.PrintDefaults()
		p("\nIf program is not given, the following will be tried:\n")
		p(defaultProgsDescription)
		p("\n")

		if writeDirsOnSignal {
			p("\nThe list of watched directories can be dumped to a file")
			p("\n%s by sending HUP (-1) signal.\n", writeDirsOutputPattern)
		}
	}

	err := flagset.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "forever:", err)
		os.Exit(2)
	}
	if helpOnly {
		flagset.SetOutput(os.Stdout)
		flagset.Usage()
		os.Exit(0)
	}

	setupDebug(c.verbose)

	if rest := flagset.Args(); len(rest) > 0 {
		c.prog.path = rest[0]
		c.prog.args = rest[1:]
		c.prog.explicit = true
	}
	return
}

func main() {
	config := parseArgs()

	err := os.Chdir(config.dir)
	if err != nil {
		fatal("cannot prepare:", err)
	}

	w, err := newWatcher(config.timeslot)
	if err != nil {
		fatal("cannot start watcher:", err)
	}

	// watcher should add all files before looping
	w.feed()
	w.installSignal()

	go loop(w, &config.prog)
	select {} // block forever
}
