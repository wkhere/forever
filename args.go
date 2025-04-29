package main

import (
	"fmt"
	"io"
	"time"
)

func usage(w io.Writer, defaults *config) {
	p := func(format string, a ...any) { fmt.Fprintf(w, format, a...) }

	p(`Usage: forever [-d dir] [-t duration] [-m duration] [-v] [-- program ...]

  -d, --dir directory       switch to directory (default %[1]q)
  -t, --delay duration      delay of write events (default %[2]v)
  -m, --min-run duration    minimal run duration (default %[3]v); if the program
                            was faster, there is a wait before further actions.
  -v, --verbose             be verbose
  -h, --help                show this help and exit
`,
		defaults.dir,
		defaults.delay,
		defaults.minRun,
	)

	p("\nIf program is not given, the following will be tried:\n")
	p(defaultProgsDescription)
	p("\n")

	if writeDirsOnSignal {
		p("\nThe list of watched directories can be dumped to a file")
		p("\n%s by sending HUP (-1) signal.\n", writeDirsOutputPattern)
	}
}

func parseArgs(args []string) (c *config, err error) {
	c = defaults()

	var rest []string

flags:
	for len(args) > 0 && err == nil {

		switch arg := args[0]; {

		case arg == "-d", arg == "--dir":
			c.dir, args, err = parseStrFlag(arg, args[1:])

		case arg == "-t", arg == "--delay":
			c.delay, args, err = parseDurationFlag(arg, args[1:])

		case arg == "-m", arg == "--min-run":
			c.minRun, args, err = parseDurationFlag(arg, args[1:])

		case arg == "-v", arg == "--verbose":
			c.verbose, args = true, args[1:]

		case arg == "-h", arg == "--help":
			c.help = func(w io.Writer) { usage(w, defaults()) }
			return c, nil

		case arg == "--":
			rest = args[1:]
			break flags

		case len(arg) > 1 && arg[0] == '-':
			return c, fmt.Errorf("unknown flag: %s", arg)

		default:
			rest = args
			break flags
		}
	}
	if err != nil {
		return c, err
	}

	setupDebug(c.verbose)

	if len(rest) > 0 {
		c.prog.path = rest[0]
		c.prog.args = rest[1:]
		c.prog.explicit = true
	}
	return c, nil
}

func parseStrFlag(s string, args []string) (x string, rest []string, err error) {
	if len(args) < 1 {
		return "", args, fmt.Errorf("flag %s: arg required", s)
	}
	return args[0], args[1:], nil
}

func parseDurationFlag(s string, args []string) (x time.Duration, rest []string, err error) {
	if len(args) < 1 {
		return 0, args, fmt.Errorf("flag %s: arg required", s)
	}
	x, err = time.ParseDuration(args[0])
	if err != nil {
		err = fmt.Errorf("flag %s: %w", s, err)
	}
	return x, args[1:], err
}
