package main

import "path/filepath"

var ignoredMounts = []string{
	"/dev",
	"/proc",
	"/sys",
	"/run",
}

var ignoredDirs = []string{
	".git",
	".stfolder",
	".stversions",
	"vendor",
	"__pycache__",
	".mypy_cache",
	"ebin",
	"deps",
	"_build",
	"classes",
}

func isInIgnoredMount(path string) bool {
	for _, p := range ignoredMounts {
		if dirContains(p, path) {
			return true
		}
	}
	return false
}

func isIgnoredDir(path string) bool {
	_, last := filepath.Split(path)
	for _, d := range ignoredDirs {
		if d == last {
			return true
		}
	}
	return false
}
