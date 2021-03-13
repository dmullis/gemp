// Copyright 2020 Donald Mullis. All rights reserved.

package internal

import (
	"bufio"
	"io"
	"log"
	"os"
)

func CountLines(r *os.File) (count int) {
	_, _ = r.Seek(0, 0)
	b := bufio.NewReader(r)
	for ; ; count++ {
		_, err := b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalf("Could not read line from \"%s\", err=%v",
				r.Name(), err)
		}
	}
}
