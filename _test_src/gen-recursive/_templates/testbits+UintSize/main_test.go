// Copyright 2020 Donald Mullis. All rights reserved.

// {{.CodeGenWarning}}

package main_test

import (
	"flag"
	"math/rand"
	"os"
	"testing"

	"github.com/dmullis/gemp/_test_src/lib"
)

type (
	baseUintType = uint{{.UintSize}}
)

var (
	numOperations = flag.Int("numOperations", 1<<10,
		`number of data elements to be operated on`)
	Data []baseUintType
)

func initData() {
	Data = make([]baseUintType, *numOperations)
	for i, _ := range Data {
		Data[i] = baseUintType(rand.Uint64())
	}
}

func TestMain(m *testing.M) {
	lib.FlagParse()
	os.Exit(m.Run())
}
