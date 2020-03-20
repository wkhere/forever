package main

import (
	"os"
	"syscall"
)

func init() {
	rusageExtras = rusageExtrasLinux{}
}

type rusageExtrasLinux struct{}

// maxRss returns RSS usage in kB.
// On Linux this is what getrusage(2) returns.
func (_ rusageExtrasLinux) maxRss(pst *os.ProcessState) (int, bool) {
	rusage, ok := pst.SysUsage().(*syscall.Rusage)
	if !ok {
		return -1, false
	}
	return int(rusage.Maxrss), true
}
