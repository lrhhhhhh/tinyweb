package tinyweb

import (
	"fmt"
	"testing"
)

/*
	1. /:a 和 /a 怎么处理？
	2. /:a 和 /:b 怎么处理? 答：每一级只能有一种冒号
*/
func TestTrie(t *testing.T) {
	paths := []string{
		"/",
		"/:lrh",
		"/:lrh/:age",
		"/a/b/c",
	}

	tree := newTrie()

	for _, x := range paths {
		tree.insert(x)
	}

	cases := []struct {
		x  string
		ok bool
	}{
		{"/lrh", true},
		{"/", true},
		{"/lrh", true},
		{"/xxxx", true},
		{"/xxxx/", true},
		{"/xxxx/xxxx", true},
		{"/xxxx/xxxx/", true},
		{"/xxxx/xxxx/xxxx", false},
		{"/a/b/c", true},
		{"/a/b/c/", true},
	}

	fmt.Println()
	for i := range cases {
		path, params, ok := tree.find(cases[i].x)
		if ok == cases[i].ok {
			fmt.Printf("found: %v\n%s ---match--> %s\npath parameters: %+v\n\n", ok, cases[i].x, path, params)
		} else {
			t.Fatalf("case %+v failed...", cases[i])
		}
	}
}
