// Copyright 2020 Donald Mullis. All rights reserved.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/dmullis/gemp/internal"
	"github.com/dmullis/gemp/internal/gen"
)

const (
	UNINITIALIZED_PATH = ""

	VALUE_LIST_COMMA_SEPARATOR = ","
	ValueListRegexp            = "[^" + VALUE_LIST_COMMA_SEPARATOR + "]+"
)

const (
	Gen  = "gen"
	Dump = "dump"
)

// General args
var (
	help = flag.Bool("h", false,
		`Repeat this message.`) // X returns status 'success' to shell
	helpAsMarkdown = flag.Bool("helpAsMarkdown", false,
		`Format help output, if any, as Markdown`)

	// X  Why not eliminate '-format' as a parameter for 'gen'?  Because
	//    'format' can introduce additional chars into the name of the output
	//    file to help set it off from the input filename.
	//
	// X  Output from 'dump' generally needs a final newline, whereas
	//   'gen' output (into names of generated files) generally does not.
	// X  flag.String() replaces any "\n" sequence it reads with "\\n", causing
	//    fmt.Printf(*format, ...) to produce wrong output.
	// XX As an expedience, 'gen' will always tack on a final newline.
	format = flag.String("format",
		"%-.s-%s", // X  'gen' -- replace '+' with '-', Key with Value
		//"%s=%s",           // X  'dump' -- for reading by 'sh'
		//"const %s=\"%s\"", // X  'dump' -- Go or JavaScript

		`Format string syntax is that of Go's 'fmt' package, with exactly
two string expansion codes e.g. "%s-%s" required.

Each pair of Key, Value strings is passed to fmt.Printf, along with this
format string.
  'gen'  Result is reinserted into each file's output pathname
  'dump' Results written line-by-line to stdout.
Note that prefixing with '%-.s', drops a string from output.
`)

	// XX  Add to test suite.
	KVplusPath = flag.String("kvpluspath", UNINITIALIZED_PATH,
		`Alternative to specifying K=V+ pairs on the command line. Arg is a
path to an input file containing Key=Value+ pairs, in 'sh' syntax.
Lines of commentary, beginning with '#', are ignored.`)

	verbose = flag.Bool("verbose", false,
		`Log heavily`)

	commandName string
	//kvpArgs []internal.KvpArg // preserves original order of keys
)

func init() {
	// X  Any environment variables specified in the form 'K=\tV' appear here as 'K=\\tV'
	// envir := os.Environ()

	log.SetFlags(log.Lshortfile)
}

func usage() {
	fpf := func(format string, s ...interface{}) {
		fmt.Fprintf(os.Stderr, format, s...)
	}
	toggleCode := func() {
		if *helpAsMarkdown {
			fmt.Fprintf(os.Stderr, "```\n")
		}
	}

	fpf("Usage:\n")
	toggleCode()

	fpf("%s", cliUsage())
	flag.PrintDefaults()

	fpf("  (K=V1,V2...Vn)*")
	fpf(`
        Any number of Key=Value+ pairs, where Value+ may be a comma-
	separated list of multiple string values to be substituted serially
	into each of multiple output directories or files.
`)
	toggleCode()
	commandSynopsis := `
 gen

  'gen' scans a single named input file in the format specified by the
  Go standard library 'template' package.  If an expansion of a known
  Key is found, each of its Values is iteratively substituted
  in, with output written to newly created files.  A K=V1,V2,...Vn pair
  multiplies the number of output files by 'n',
  with successive files receiving V1,V2...Vn for substitution within
  the file.
  For 'gen'-specific help:
      $ gemp _=_ gen

 dump

  'dump' reads a single-line format string from the command line.
  Result is written to stdout, with the format string expanded by each
  Key-Value pair on successive lines of the output.  Any value list
  V1,V2...Vn passed to 'dump' is not expanded or parsed further but
  merely treated as a single string.
`
	fpf("%s", commandSynopsis)
}

func cliUsage() string {
	exeName := path.Base(os.Args[0])
	var flagUsage string
	flag.VisitAll(
		func(f *flag.Flag) {
			flagUsage += fmt.Sprintf("[-%s=%s] ", f.Name, f.DefValue)
		})
	return fmt.Sprintf("%s %s[K=V1,V2...Vn]* (gen [flags] input_file | dump)\n\n",
		exeName, flagUsage)
}

func main() {
	flag.Usage = usage

	flag.Parse()
	if !flag.Parsed() {
		log.Fatalln("flag.Parsed() == false")
	}

	kvpArgs, nonKvpArgs := getKVplus()
	if len(nonKvpArgs) == 0 {
		if *help {
			usage()
			os.Exit(0)
		} else {
			usageWhy("No command found")
		}
	}
	commandName := nonKvpArgs[0]

	if *help {
		switch commandName {
		case Gen:
			gen.UsageDump(*helpAsMarkdown, cliUsage())
		case Dump:
			usage()
		}
		os.Exit(0)
	}

	if len(kvpArgs) < 1 {
		usageWhy("\nno Key=Value+ pairs found")
	}

	switch commandName {
	case Gen:
		numLines := gen.ParseArgs(nonKvpArgs[1:], cliUsage())
		gen.ExpandTemplate(*verbose, *format, numLines,
			kvpArgs)
	case Dump:
		for _, kvp := range kvpArgs {
			kvpValues := kvp.Values[0]
			for _, v := range kvp.Values[1:] {
				kvpValues += VALUE_LIST_COMMA_SEPARATOR + v
			}
			expand := fmt.Sprintf(*format, kvp.Key, kvpValues)
			if expand[:2] == "%!" {
				usageWhy(fmt.Sprintf(
					"format (%s) failed to expand arg (%s=%s):\n\t%s\n",
					*format, kvp.Key, kvpValues, expand))
			}
			if _, err := io.WriteString(os.Stdout, expand+"\n"); err != nil {
				usageWhy(fmt.Sprintf("io.WriteString() failed: %v\n", err))
			}
		}
		return
	default:
		usageWhy(fmt.Sprintf("No such command: %s", commandName))
	}
}

func usageWhy(why string) {
	usage()
	fmt.Fprintf(os.Stderr, "\n%s\n\n", why)
	os.Exit(1)
}

func getKVplus() (kvpArgs []internal.KvpArg, remainingArgs []string) {
	kvpArgs, remainingArgs = scanForKVplusArgs(flag.Args())

	// XX  File must be parsed before KV pairs on command line if latter are to override.
	if *KVplusPath != "" {
		kvpArgs = append(kvpArgs, scanKVplusFile(*KVplusPath)...)
	}
	return
}

func scanForKVplusArgs(args []string) (
	kvpArgs []internal.KvpArg, remainingArgs []string) {

	for iArg, arg := range args {
		kvp := strings.Split(arg, "=")
		if len(kvp) != 2 {
			remainingArgs = args[iArg:]
			break
		}
		newKvpArg := newKVplusPair(kvp)

		// Search earlier Keys for duplicates.
		//   XX  N^2 in number of Keys -- use a map instead?
		for _, kvp := range kvpArgs {
			if kvp.Key == newKvpArg.Key {
				// XX  ? Add option to accumulate the values K=V1, K=V2, ... ,
				//       as an alternative to comma-separated syntax K=V1,V2,...
				log.Fatalf("Duplicate key specified: '%v', '%v'", kvp, newKvpArg)
			}
		}
		kvpArgs = append(kvpArgs, newKvpArg)
	}
	return
}

func newKVplusPair(newKvp []string) internal.KvpArg {
	vetKVstring(newKvp)
	return parseKvpArg(newKvp)
}

func vetKVstring(kvplus []string) {
	reportFatal := func(format string) {
		caller := func(howHigh int) string {
			pc, file, line, ok := runtime.Caller(howHigh)
			_ = pc
			if !ok {
				return ""
			}
			baseFileName := file[strings.LastIndex(file, "/")+1:]
			return baseFileName + ":" + strconv.Itoa(line)
		}
		fmt.Fprintf(os.Stderr, format, kvplus)
		//fmt.Fprintf(os.Stderr, "Reading template file %s", templatePath)
		fmt.Fprintf(os.Stderr, " Called from line %3s <- %3s <- %3s\n",
			caller(2), caller(3), caller(4))
		//usage()
		log.Fatalln("FATAL")
	}
	if len(kvplus) != 2 {
		reportFatal("Appears not to be a Key=Value+ pair: %v\n")
	}
	if len(kvplus[0]) <= 0 {
		reportFatal("Key side of Key=Value+ pair empty: \"%v\"\n")
	}
	if len(kvplus[1]) <= 0 {
		reportFatal("Value+ side of Key=Value+ pair empty: \"%v\"\n")
	}
}

func parseKvpArg(rawKvp []string) (kvpArg internal.KvpArg) {
	kvpArg.Key = rawKvp[0]
	valueStringsRE := regexp.MustCompile("(" + ValueListRegexp + ")")
	commaSeparatedValues := valueStringsRE.FindAllStringSubmatch(rawKvp[1], -1)
	if len(commaSeparatedValues) < 1 {
		log.Fatalf("Key= '%s=' specified, but no value found on RHS", kvpArg.Key)
	}
	// split out the Values into a slice of string
	for _, match := range commaSeparatedValues {
		kvpArg.Values = append(kvpArg.Values, match[1])
	}
	return
}

func scanKVplusFile(kVplusPath string) (kvpArgs []internal.KvpArg) {
	kvfile, err := os.Open(kVplusPath)
	if err != nil {
		log.Fatalln(err)
	}
	scanner := bufio.NewScanner(kvfile)
	// Iterate over each (non-comment) line of file contents.
	for scanner.Scan() {
		kvp := strings.Split(scanner.Text(), "=")
		kvp[0] = strings.TrimSpace(kvp[0])
		if len(kvp[0]) == 0 || kvp[0][0] == '#' {
			continue
		}
		if len(kvp) < 2 {
			log.Fatalf("len(kvp)==%d < 2; kvp==%v", len(kvp), kvp)
		}
		// find end of RHS 'Value+' string
		if len(kvp[1]) >= 2 && strings.IndexAny(string(kvp[1][0]), "'\"") >= 0 {
			quoteChar := kvp[1][0]
			close := strings.IndexByte(string(kvp[1][1:]), quoteChar)
			kvp[1] = kvp[1][1 : 1+close]
		} else { // no quoting found, strip from leftmost white space onward
			close := strings.IndexAny(string(kvp[1]), " \t")
			if close >= 0 {
				kvp[1] = kvp[1][:close]
			}
		}
		kvpArgs = append(kvpArgs, newKVplusPair(kvp))
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	return
}
