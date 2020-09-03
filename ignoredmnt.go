package main

import "path/filepath"

var ignoredMounts = []string{
	"/dev",
	"/proc",
	"/sys",
	"/run",
}

func isInIgnoredMount(path string) bool {
	for _, p := range ignoredMounts {
		if dirContains(p, path) {
			return true
		}
	}
	return false
}

func dirContains(base, path string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	if rel == ".." || (len(rel) >= 3 && rel[:3] == "../") {
		return false
	}
	return true
}
