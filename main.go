package main

import (
	"fmt"
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

	for _, sd := range sds {
		diffURL := filepath.Join(dataDir, sd.File)
		funcs := ParseFile(diffURL)
		f := funcToMap(funcs)

		for _,  sdFunc := range sd.Funcs {
			if loc, ok := f[sdFunc]; ok {
				fmt.Printf("%s %s %d\n", sd.File, sdFunc, loc)
			} else {
				fmt.Println("Cannot find it!")
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
	dataDir := "./data"
	diffFile := "./data/openssl-1.0.2k-no-ssl2.patch"
	Unify(dataDir, diffFile)
	//diffDemo()
}
