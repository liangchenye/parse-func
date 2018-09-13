package main

import (
	"fmt"
	"io/ioutil"
)

func diffDemo() {
	diffFile := "./data/openssl-1.0.2k-req-x509.patch"

        diffData, err := ioutil.ReadFile(diffFile)
        if err != nil {
                return 
        }

        data := ParseData(diffData)
        items, _ := NewDiffItems(data, SearchCondition{})

	fmt.Println(items)

}

func functionDemo() {
	opensslDir := "./data/openssl-1.0.2k/"

	dirs := WalkDir(opensslDir)

	for _, d := range dirs {
		fmt.Println(d)
		funcs := ParseFile(d)
		fmt.Printf("%s has %d funcs\n", d, len(funcs), funcs)
	}
}

func main() {
	functionDemo()
}
