package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	funcs := ParseFile("testdata/test001.data")
	if len(funcs) != 1 {
		fmt.Println("invalid items ", funcs)
		return
	}

	if funcs[0].Name != "abcd" || funcs[0].LOC != 3 {
		fmt.Printf("invalid content <%s>, <%d>\n", funcs[0].Name, funcs[0].LOC)
	}
}

func TestParseFunction(t *testing.T) {
	data := `
ab int abcd()
{
	hello 
}

{
dasf
}


sdf string kb()
`
	lines := strings.Split(data, "\n")
	funcs := ParseFunction(lines)

	if len(funcs) != 1 {
		fmt.Println("invalid items ", funcs)
	}

	if funcs[0].Name != "abcd" || funcs[0].LOC != 3 {
		fmt.Printf("invalid content <%s>, <%d>\n", funcs[0].Name, funcs[0].LOC)
	}
}

func TestCheckFunction(t *testing.T) {
	var ret string
	ret = CheckFunction("abcd()")
	if ret != "abcd" {
		fmt.Println("error: ", ret)
	}

	ret = CheckFunction("abcd ()")
	if ret != "abcd" {
		fmt.Println("error: ", ret)
	}

	ret = CheckFunction("abcd();")
	if ret != "" {
		fmt.Println("error: ", ret)
	}

	ret = CheckFunction("abcd(")
	if ret != "abcd" {
		fmt.Println("error: ", ret)
	}

	ret = CheckFunction("abcd = ab()")
	if ret != "" {
		fmt.Println("error: ", ret)
	}

	ret = CheckFunction(" * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)")
	if ret != "" {
		fmt.Println("error: ", ret)
	}

	ret = CheckFunction("int app_RAND_load_file(const char *file, BIO *bio_e, int dont_warn)")
	if ret != "app_RAND_load_file" {
		fmt.Println("error: ", ret)
	}
}
