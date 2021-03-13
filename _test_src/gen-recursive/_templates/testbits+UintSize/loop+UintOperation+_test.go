// Copyright 2020 Donald Mullis. All rights reserved.

// {{.CodeGenWarning}}

// For a different way to apply the 'template' library to this class of problem:
//     - /usr/local/go/src/cmd/compile/internal/gc/testdata/gen/arithConstGen.go
//     - /usr/local/go/src/cmd/compile/internal/gc/testdata/gen/arithBoundaryGen.go

package main_test

import (
	"math/bits"
	"testing"
)

func Benchmark{{.UintOperation}}_{{.UintSize}}(b *testing.B) {
	initData()
	for i:=0; i<b.N; i++ {
		for n, _ := range Data {
			Data[n] = bits.{{.UintOperation}}{{.UintSize}}(Data[n])
		}
	}
}
