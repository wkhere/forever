package main

import (
	"testing"
)

func TestDirContains(t *testing.T) {
	tcases := []struct {
		name    string
		base, p string
		want    bool
	}{
		{"zeroall", "", "", true},
		{"zerobase absdir", "", "/foo", false},
		{"zeropath", "/foo", "", false},
		{"absdir self", "/dev", "/dev", true},
		{"absdir-slash self", "/dev/", "/dev", true},
		{"absdir ok", "/dev", "/dev/foo", true},
		{"absdir-slash ok", "/dev/", "/dev/foo", true},
		{"absdir substr", "/dev", "/devfoo", false},
		{"absdir-slash ok", "/dev/", "/dev/foo", true},
		{"illegal dot", "/dev", ".", false},
	}

	for _, tc := range tcases {
		have := dirContains(tc.base, tc.p)
		if have != tc.want {
			t.Errorf("%s: have %v, want %v", tc.name, have, tc.want)
		}
	}
}
