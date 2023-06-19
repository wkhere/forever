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
func (rusageExtrasLinux) getMemStats(pst *os.ProcessState) (memStats, bool) {
	rusage, ok := pst.SysUsage().(*syscall.Rusage)
	if !ok {
		return memStats{}, false
	}
	return memStats{
		maxRss: rusage.Maxrss,
		minFlt: rusage.Minflt,
		majFlt: rusage.Majflt,
	}, true
}
