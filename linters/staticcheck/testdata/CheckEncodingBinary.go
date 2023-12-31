package pkg

import (
	"encoding/binary"
	"io/ioutil"

	"github.com/hellobchain/wswlog/wlogging"
)

var logger = wlogging.MustGetLoggerWithoutName()

func fn() {
	type T1 struct {
		A int32
	}
	type T2 struct {
		A int32
		B int
	}
	type T3 struct {
		A []int32
	}
	type T4 struct {
		A *int32
	}
	type T5 struct {
		A int32
	}
	type T6 []byte

	var x1 int
	var x2 int32
	var x3 []int
	var x4 []int32
	var x5 [1]int
	var x6 [1]int32
	var x7 T1
	var x8 T2
	var x9 T3
	var x10 T4
	var x11 = &T5{}
	var x13 []byte
	var x14 *[]byte
	var x15 T6
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x1)) // MATCH /cannot be used with binary.Write/
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x2))
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x3)) // MATCH /cannot be used with binary.Write/
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x4))
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x5)) // MATCH /cannot be used with binary.Write/
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x6))
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x7))
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x8))  // MATCH /cannot be used with binary.Write/
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x9))  // MATCH /cannot be used with binary.Write/
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x10)) // MATCH /cannot be used with binary.Write/
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x11))
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, &x13))
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, &x14)) // MATCH /cannot be used with binary.Write/
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, x15))
	logger.Info(binary.Write(ioutil.Discard, binary.LittleEndian, &x15))
}
