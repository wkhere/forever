package main

import "os"

type rusageExtrasReader interface {
	// maxRss returns RSS usage in kBytes, plus a bool flag if it was present.
	maxRss(*os.ProcessState) (int, bool)
}

var rusageExtras rusageExtrasReader = noRusageExtras{}

type noRusageExtras struct{}

func (_ noRusageExtras) maxRss(*os.ProcessState) (int, bool) {
	return -1, false
}
