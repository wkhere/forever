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

	help func()
}

func parseArgs(args []string) (c *config, _ error) {
	c = new(config)

	var help bool
	flagset := pflag.NewFlagSet("flags", pflag.ContinueOnError)
	flagset.SortFlags = false

	flagset.StringVarP(&c.dir, "dir", "d", ".",
		"switch to `directory`")
	flagset.DurationVarP(&c.timeslot, "timeslot", "t", 200*time.Millisecond,
		"timeslot for write events")
	flagset.BoolVarP(&c.verbose, "verbose", "v", false,
		"be verbose")
	flagset.BoolVarP(&help, "help", "h", false,
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

	err := flagset.Parse(args)
	if err != nil {
		return nil, err
	}
	if help {
		c.help = func() {
			flagset.SetOutput(os.Stdout)
			flagset.Usage()
		}
	}

	setupDebug(c.verbose)

	if rest := flagset.Args(); len(rest) > 0 {
		c.prog.path = rest[0]
		c.prog.args = rest[1:]
		c.prog.explicit = true
	}
	return c, nil
}

func main() {
	config, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}
	if config.help != nil {
		config.help()
		os.Exit(0)
	}

	err = run(config)
	if err != nil {
		die(1, err)
	}
}

func die(exitcode int, err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(exitcode)
}

func run(c *config) error {

	err := os.Chdir(c.dir)
	if err != nil {
		return err
	}

	w, err := newWatcher(c.timeslot)
	if err != nil {
		return err
	}

	// watcher should add all files before looping
	err = w.feed()
	if err != nil {
		return err
	}
	w.installSignal()

	go loop(w, &c.prog)
	select {} // block forever
}
