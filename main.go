package main

import "fmt"

func main() {
	content := "abcd int ab(dadf)"
	funcStr := CheckFunction(content)
	fmt.Println(funcStr)
}

