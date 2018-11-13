// +build linux darwin solaris dragonfly freebsd netbsd openbsd

package main

import (
	"os"
	"syscall"
)

func init() {
	sysRusageExtras = rusageExtrasUnix{}
}

type rusageExtrasUnix struct{}

func (x rusageExtrasUnix) maxRss(pst *os.ProcessState) (int64, bool) {
	rusage, ok := pst.SysUsage().(*syscall.Rusage)
	if !ok {
		return -1, false
	}
	return rusage.Maxrss, true
}
