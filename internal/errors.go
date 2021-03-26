// Copyright 2020 Donald Mullis. All rights reserved.

package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const MarkdownAutoGenMessage = "<!-- DO NOT MODIFY -- automatically generated -->"

func ToggleCode(header string) {
	fmt.Fprintf(os.Stderr, "%s\n```\n", header)
}

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
