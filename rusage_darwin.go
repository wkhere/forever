package main

import (
	"os"
	"syscall"
)

func init() {
	rusageExtras = rusageExtrasDarwin{}
}

type rusageExtrasDarwin struct{}

// maxRss returns RSS usage in kB.
// On Darwin getrusage(2) returns bytes count.
func (_ rusageExtrasDarwin) maxRss(pst *os.ProcessState) (int, bool) {
	rusage, ok := pst.SysUsage().(*syscall.Rusage)
	if !ok {
		return -1, false
	}
	return int(rusage.Maxrss / 1024), true
}
