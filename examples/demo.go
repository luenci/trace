package main

func foo() {
	bar()
}
func bar() {
}
func demo() {
	foo()
}

func main() {
	demo()
}

/*
out:

[./bin/core examples/demo.go]
package main

import "github.com/luenci/trace"

func foo() {
	defer trace.Trace()()
	bar()
}
func bar() {
	defer trace.Trace()()
}
func demo() {
	defer trace.Trace()()
	foo()
}

func main() {
	defer trace.Trace()()
	demo()
}
*/
