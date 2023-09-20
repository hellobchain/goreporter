package pkg

import (
	"fmt"

	"github.com/hellobchain/wswlog/wlogging"
)

var logger = wlogging.MustGetLoggerWithoutName()

func fn() {
	var s string
	fn2 := func() string { return "" }
	fmt.Printf(fn2())   // MATCH /should use print-style function/
	fmt.Sprintf(fn2())  // MATCH /should use print-style function/
	logger.Infof(fn2()) // MATCH /should use print-style function/
	fmt.Printf(s)       // MATCH /should use print-style function/
	fmt.Printf(s, "")

	fmt.Printf(fn2(), "")
	fmt.Printf("")
	fmt.Printf("", "")
}
