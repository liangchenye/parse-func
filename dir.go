package main

import (
	"path/filepath"
	"os"
	"strings"
)

func WalkDir(topDir string) []string{
	var files []string
	filepath.Walk(topDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".c") {
			files = append(files, path)
		}
		return nil
	})

	return files
}
