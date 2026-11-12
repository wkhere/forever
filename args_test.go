package main

import (
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"mvdan.cc/sh/v3/shell"
)

func TestParseArgs(t *testing.T) {
	var (
		dir = func(s string) func(*config) bool {
			return func(c *config) bool { return c.dir == s }
		}
		single = func(x bool) func(*config) bool {
			return func(c *config) bool { return c.single == x }
		}
		delay = func(x time.Duration) func(*config) bool {
			return func(c *config) bool { return c.delay == x }
		}
		minRun = func(x time.Duration) func(*config) bool {
			return func(c *config) bool { return c.minRun == x }
		}
		v = func(x bool) func(*config) bool {
			return func(c *config) bool { return c.verbose == x }
		}
		progPath = func(p string) func(*config) bool {
			return func(c *config) bool {
				return c.prog.path == p
			}
		}
		progArgs = func(a ...string) func(*config) bool {
			return func(c *config) bool {
				return slices.Equal(c.prog.args, a)
			}
		}
		help = func(pattern string) func(*config) bool {
			return func(c *config) bool {
				if c.help == nil {
					return false
				}
				b := new(strings.Builder)
				c.help(b)
				return strings.Contains(b.String(), pattern)
			}
		}
		none = func(*config) bool { return false }
		all  = func(ff ...func(*config) bool) func(*config) bool {
			return func(c *config) bool {
				for _, f := range ff {
					if !f(c) {
						return false
					}
				}
				return true
			}
		}
		d        = defaults()
		defaults = all(dir(d.dir), delay(d.delay), progPath(""), progArgs())
	)

	var tab = []struct {
		input string
		errs  string
		want  func(*config) bool
	}{
		{"", "", defaults},
		{"--", "", defaults},
		{"-h", "", help("Usage:")},
		{"--help", "", help("Usage:")},
		{"-q", "unknown flag -q", none},
		{"-q -x", "unknown flag -q", none},
		{"--quux", "unknown flag --quux", none},
		{"-q -h", "unknown flag -q", none},
		{"-q --help", "unknown flag -q", none},
		{"-h -q", "", help("Usage:")},
		{"--help -q", "", help("Usage:")},
		{"./nonexistent", "", progPath("./nonexistent")}, // prog gets checked after parseArgs
		{"ls", "", all(progPath("ls"), v(false))},
		{"ls 1 2", "", all(progPath("ls"), progArgs("1", "2"), v(false))},
		{"-d", "needs a value", none},
		{"-d /tmp", "", all(dir("/tmp"))},
		{"--dir", "needs a value", none},
		{"--dir /tmp", "", all(dir("/tmp"))},
		{"-t", "needs a value", none},
		{"-t foo", "invalid duration", none},
		{"-t 700ms", "", delay(700 * time.Millisecond)},
		{"--delay", "needs a value", none},
		{"--delay foo", "invalid duration", none},
		{"--delay 700ms", "", delay(700 * time.Millisecond)},
		{"-m", "needs a value", none},
		{"-m foo", "invalid duration", none},
		{"-m 300ms", "", minRun(300 * time.Millisecond)},
		{"--min-run", "needs a value", none},
		{"--min-run foo", "invalid duration", none},
		{"--min-run 300ms", "", minRun(300 * time.Millisecond)},
		{"-1 ls", "", all(single(true), progPath("ls"), progArgs(), v(false))},
		{"-1 -- ls", "", all(single(true), progPath("ls"), progArgs(), v(false))},
		{"-v ls", "", all(v(true), progPath("ls"), progArgs())},
		{"--verbose ls", "", all(v(true), progPath("ls"), progArgs())},
		{"-v -1 ls", "", all(single(true), v(true), progPath("ls"), progArgs())},
		{"-v ls -1", "", all(single(false), v(true), progPath("ls"), progArgs("-1"))},
		{"-v -- ls -1", "", all(single(false), v(true), progPath("ls"), progArgs("-1"))},
	}

	for i, tc := range tab {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			args, err := shell.Fields(tc.input, nil)
			if err != nil {
				t.Fatalf("shell args: %v", err)
			}
			c, err := parseArgs(args)

			switch {
			case err != nil && tc.errs == "":
				t.Errorf("unexpected error: %s", err)

			case err != nil && tc.errs != "":
				if !strings.Contains(err.Error(), tc.errs) {
					t.Errorf("expected error with %q; having: %q", tc.errs, err)
				}

			case err == nil && tc.errs != "":
				t.Errorf("expected error with %q", tc.errs)

			default:
				if !(tc.want(c)) {
					t.Errorf("predicate failed on %+v", c)
				}
			}
		})
	}
}
