// Test of return+else warning; should not trigger on multi-branch if/else.
// OK

// Package pkg ...
package pkg

import (
	"github.com/hellobchain/wswlog/wlogging"
)

var logger = wlogging.MustGetLoggerWithoutName()

func f(x int) bool {
	if x == 0 {
		logger.Info("x is zero")
	} else if x > 0 {
		return true
	} else {
		logger.Infof("non-positive x: %d", x)
	}
	return false
}
