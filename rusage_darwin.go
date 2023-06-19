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
func (_ rusageExtrasDarwin) getMemStats(pst *os.ProcessState) (memStats, bool) {
	rusage, ok := pst.SysUsage().(*syscall.Rusage)
	if !ok {
		return memStats{}, false
	}
	return memStats{
		maxRss: rusage.Maxrss / 1024,
		minFlt: rusage.Minflt,
		majFlt: rusage.Majflt,
	}, true
}
