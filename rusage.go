package main

import "os"

type rusageExtras interface {
	maxRss(*os.ProcessState) (int64, bool)
}

var sysRusageExtras rusageExtras = noRusageExtras{}

type noRusageExtras struct{}

func (_ noRusageExtras) maxRss(*os.ProcessState) (int64, bool) {
	return -1, false
}
