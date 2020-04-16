package g_util

import (
	"fmt"
	"testing"
)

var c = make(chan bool)

func TestDefer(t *testing.T) {
	go run()
	_ = <-c
	fmt.Println("test case finished.")
}

func run() {
	defer func() {
		fmt.Println("exit func run")
	}()
	run2()
}

func run2() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic catch %s\n", err)
		}
		fmt.Println("exit func run2")
		close(c)
	}()
	panic("abc")
}
