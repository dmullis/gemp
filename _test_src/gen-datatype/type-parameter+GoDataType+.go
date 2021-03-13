// {{.Copyright}}{{- /*
//    pad for line 2 of copyright notice
//    pad for line 3 of copyright notice*/}}

// {{if false}}
// +build ignore
// {{else}} {{.CodeGenWarning}} {{.NewLine}} {{.NewLine}} {{end}}

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var zeroValue {{.GoDataType}}

	fmt.Printf("Zero value of type '{{.GoDataType}}' is \"%#v\"\n", zeroValue)
	fmt.Printf("Runtime reports type '%T' consumes %d bytes\n",
		zeroValue, int(unsafe.Sizeof(zeroValue)))
}
