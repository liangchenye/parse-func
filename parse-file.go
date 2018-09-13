package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Function struct {
	Name string
	LOC int
}

func init() {
	fmt.Println("hello, wolrd!")
}

func ParseFile(filename string) []Function {
	var funcs []Function
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return funcs
	}

	lines := strings.Split(string(data), "\n")
	return ParseFunction(lines)
}

func ParseFunction(lines []string) []Function{
	fBegin := ""
	lBegin := false

	count := 0
	var funcs []Function

	for i:=0; i < len(lines); i++ {
		if lines[i] == "{" {
			if fBegin != "" {
				lBegin = true
				count ++
			}
			continue
		}

		if lBegin {
			count ++

			if lines[i] == "}" {
				var fu Function
				fu.Name = fBegin 
				fu.LOC = count
				funcs = append(funcs, fu)

				fBegin = ""
				lBegin = false
				count = 0
			}
		} else {
			f := CheckFunction(lines[i])
			if f!= "" {
				fBegin = f
			}
		}
	}

	return funcs
}

func CheckFunction(line string) string {
	begin := -1
	end := -1
	left := false
	right := false

	constStr := " ()"
	constBlank := 0
	constLeft := 1
	constRight := 2

	invalidStr := "{}#=;"
	for i:= 0; i < len(line); i++ {
		switch  line[i] {
		case constStr[constBlank] :
			// It might be a bug.. not sure
			// ignore inside (...)
			if left == false {
				begin = -1
			}
		case constStr[constLeft] :
			if begin > -1 {
				end = i
				left = true
			}
			right = false
		case  constStr[constRight] :
			if left {
				right = true
			}
		case  invalidStr[0]:
			return ""
		case  invalidStr[1]:
			return ""
		case  invalidStr[2]:
			return ""
		case  invalidStr[3]:
			return ""
		case  invalidStr[4]:
			return ""
		default:
			if begin > -1 {
				right = false
			} else {
				begin = i
			}
		}
	}

	if left {
// we might not always get ')'
// int abcd(int a
//          int b)
//{
//}
//	if left && right {
		return line[begin:end]
		fmt.Println(right)
	}

	return ""
}
