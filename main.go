package main

import "fmt"

func main() {
	opensslDir := "./data/openssl-1.0.2k/"

	dirs := WalkDir(opensslDir)
//	fmt.Println(dirs)

	for _, d := range dirs {
		fmt.Println(d)
		funcs := ParseFile(d)
		fmt.Printf("%s has %d funcs\n", d, len(funcs), funcs)
		break
	}
}

