package main

import (
	"fmt"
	"testing"
)

func TestWalkDir(t *testing.T) {
	paths := WalkDir("testdata")

	if len(paths) != 1 {
		fmt.Println("Invalid struct")
		return
	}

	if paths[0] != "testdata/test001.c" {
		fmt.Println("Invalid filename")
	}
}
