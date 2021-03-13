// Copyright 2020 Donald Mullis. All rights reserved.

package lib

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func FlagParse() {
	usageWriter := flag.CommandLine.Output()

	usageFunc := func() {
		fmt.Fprintf(usageWriter, "Usage of %s:\n", os.Args[0])
	}
	flag.Usage = func() {
		usageFunc()
		flag.PrintDefaults()
		printFlags(os.Stderr)
		os.Exit(0)
	}
	fail := func(err error) {
		fmt.Fprintf(usageWriter, "%v\n\n", err)
		usageFunc()
		os.Exit(1)

	}
	if len(os.Args) <= 1 {
		fail(errors.New("no command line arguments"))
	}
	if len(flag.Args()) > 0 {
		fail(errors.New("garbage on command line: " + flag.Args()[0]))
	}

	flag.Parse()
	if !flag.Parsed() {
		log.Fatalln("! flag.Parsed()")
	}
}

func printFlags(outFd *os.File) {
	flag.Visit(
		func(pf *flag.Flag) {
			outFormat := "\t%-32s %s\n"
			if len(pf.Value.String()) == 0 {
				fmt.Fprintf(outFd, outFormat, pf.Name, "-")
				return
			}
			fmt.Fprintf(outFd, outFormat, pf.Name, pf.Value)
		})
}
