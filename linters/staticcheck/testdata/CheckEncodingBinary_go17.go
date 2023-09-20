package pkg

import (
	"encoding/binary"
	"io/ioutil"

	"github.com/hellobchain/wswlog/wlogging"
)

var logger = wlogging.MustGetLoggerWithoutName()

func fn() {
	var x bool
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x)) // MATCH "cannot be used with binary.Write"
}
