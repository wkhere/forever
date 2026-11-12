package main

import (
	"fmt"
	"io"

	"codeberg.org/wkhere/argp"
	pt "codeberg.org/wkhere/argp/parsetime"
)

func usage(w io.Writer) {
	p := func(format string, a ...any) { fmt.Fprintf(w, format, a...) }
	d := defaults()

	p(`Usage: forever [-1] [-d dir] [-m duration] [-t duration] [-v] [-- program ...]

  -1, --single              restrict watch to just one non-recursive directory
  -d, --dir directory       switch to directory (default %[1]q)
  -m, --min-run duration    minimal run duration (default %[2]v); if the program
                            was faster, there is a wait before further actions.
  -t, --delay duration      delay of write events (default %[3]v)
  -v, --verbose             be verbose
  -h, --help                show this help and exit
`,
		d.dir,
		d.minRun,
		d.delay,
	)

	p("\nIf program is not given, the following will be tried:\n")
	p(defaultProgsDescription)

	if writeDirsOnSignal {
		p("\nThe list of watched directories can be dumped to the file %s\n", writeDirsOutputPattern)
		p("by sending HUP (-1) signal.\n")
	}
}

func parseArgs(args []string) (c *config, _ error) {
	c = defaults()

	var p = &argp.PState{Args: args}
	var rest []string

flags:
	for ; len(p.Args) > 0 && p.Err == nil; p.Next() {
		switch {

		case p.MatchFlagExpr("-d", "--dir"):
			c.dir = p.Val

		case p.MatchJustFlag("-1", "--single"):
			c.single = true

		case p.MatchFlagExpr("-t", "--delay"):
			pt.ParseDuration(&c.delay, p)

		case p.MatchFlagExpr("-m", "--min-run"):
			pt.ParseDuration(&c.minRun, p)

		case p.MatchJustFlag("-v", "--verbose"):
			c.verbose = true

		case p.MatchJustFlag("-h", "--help"):
			c.help = usage
			return c, nil

		case p.Arg() == "--":
			rest = p.Args[1:]
			break flags

		case p.MatchAnyFlag():
			p.Errorf("unknown flag %s", p.Flag)

		default:
			rest = p.Args
			break flags
		}
	}

	if p.Err != nil {
		return c, p.Err
	}

	setupDebug(c.verbose)

	if len(rest) > 0 {
		c.prog.path = rest[0]
		c.prog.args = rest[1:]
		c.prog.explicit = true
	}
	return c, nil
}
