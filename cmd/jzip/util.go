package main

import (
	"os"
	"path/filepath"
)

func absolutePath(path string) string {
	if path != "" && path[0] != '/' {
		wd, err := os.Getwd()
		if err == nil {
			path = filepath.Join(wd, path)
		}
	}
	return path
}
