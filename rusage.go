package main

import "os"

type rusageExtras interface {
	maxRss(*os.ProcessState) (int, bool)
}

var sysRusageExtras rusageExtras = noRusageExtras{}

type noRusageExtras struct{}

func (_ noRusageExtras) maxRss(*os.ProcessState) (int, bool) {
	return -1, false
}
