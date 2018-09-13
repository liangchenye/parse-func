package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"path/filepath"
)

func funcToMap(funcs []Function) map[string]int {
	m := make(map[string]int)
	for _, f := range funcs {
		m[f.Name] = f.LOC
	}
	return m
}

func Unify(dataDir string, diffFile string) {
	sds := GetSimpleDiff(diffFile)

	fmt.Printf("Parsing %s\n", diffFile)
	for _, sd := range sds {
		diffURL := filepath.Join(dataDir, sd.File)
		funcs := ParseFile(diffURL)
		f := funcToMap(funcs)

		sdMaps := make(map[string]bool)
		for _,  sdFunc := range sd.Funcs {
			sdMaps[sdFunc] = true
		}
		for sdFunc, _ := range sdMaps {
			if loc, ok := f[sdFunc]; ok {
				fmt.Printf("\tLiang Hit it! %s %s %d\n", sd.File, sdFunc, loc)
			} else {
				fmt.Printf("\tCannot find the function! %s %s\n", sd.File, sdFunc)
			}
		}
//		fmt.Printf("%s has %d funcs\n", diffURL, len(funcs), funcs)
	}
}

type SimpleDiff struct {
	File string
	Funcs []string
}

func GetSimpleDiff(diffFile string) []SimpleDiff {
//	diffFile := "./data/openssl-1.0.2k-no-ssl2.patch"
	var sds []SimpleDiff
        diffData, err := ioutil.ReadFile(diffFile)
        if err != nil {
		fmt.Println("Fatal error!\n")
		panic("cannot get simple diff")
                return sds
        }

        data := ParseData(diffData)
        items, _ := NewDiffItems(data, SearchCondition{})

	for _, item := range items {
		origin := item.GetFile()
		if !strings.HasSuffix(origin, ".c") {
			continue
		}

		var funcs []string
		for _, frag := range item.Frags{
			ff := frag.GetFunction()
			if ff != "" {
				funcs = append(funcs, ff)
			}
		}
		if len(funcs) > 0 {
			var sd SimpleDiff
			sd.File = origin
			sd.Funcs = funcs
			sds = append(sds, sd)
		}
	}

	return sds
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
	if len(os.Args) < 3 {
		fmt.Println("help:  datadir, diffurl")
	}

	dataDir := "./data"
	diffFile := "./data/openssl-1.0.2k-no-ssl2.patch"

	dataDir = os.Args[1]
	diffFile = os.Args[2]
	Unify(dataDir, diffFile)
	//diffDemo()
}
