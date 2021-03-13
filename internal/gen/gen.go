// Copyright 2020 Donald Mullis. All rights reserved.

// Gemp 'gen' recursively expands template files as directed by
// a common pool of K=V bindings.
//
//    Tenets:
//     1. Don't break runtime's understanding of source code line numbers -- strings
//        inserted into source code by generation must not include line breaks.
package gen

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/dmullis/gemp/internal"
)

// Args specific to "gen"
var (
	usagePreamble = `command 'gen' usage:

  If a list of more than one value has been assigned to a variable 'K', 'K'
  must be expanded by the template file in order to avoid identical
  duplicate output files.

  In order to generate unique names for each output file, the Key
  introducing K=V1,V2...Vn must be made available for substitution in
  the name of the input file by setting off its name K with a
  separator character, reserved for no other use.

  K's expansion for pathnames is controlled by the general '-format='
  argument.

  Elements of 'templatePath' will be split into substrings at each
  transition from a character legal in Go identifiers '[a-zA-Z0-9_]',
  to one that is not.  Each such substring will then be tested against
  all Keys specified.  For the first matching key only, each
  of its one or more specified values will be substituted in
  turn, with corresponding separate output sub-directories and
  the base file written.  Format of generated pathnames is
  controlled by the 'format' option.

  Directory names with initial '_' are useful to hide source for code
  generation from any run of "go mod tidy" initiated at the root directory.
`
	fs = flag.NewFlagSet("gen", flag.ExitOnError)

	clobber = fs.Bool("clobber", false,
		`Overwrite already-existing output files.`)

	//   https://golang.org/pkg/path/
	//   https://golang.org/pkg/text/template/#hdr-Arguments
	inKeySeparator = fs.String("inkeyseparator",
		"", // XX  no default
		`Input files may be visually distinguished from output
files they generate by inclusion of a specified character.  The character
must not be legal in a Go identifer ([a-zA-Z0-9_]).  Any instances of
the character will be omitted from output file names.

Candidates for '-inkeyseparator' usage must seek a compromise:
   a. Escape special treatment by build tools, command shells,
      or GNU's 'readline' library, and
   b. Not collide with other non-alphanums wanted within filenames.
A few non-alphanumeric candidates: + ~ @  %`)

	outTopDir = fs.String("outtopdir", ".",
		`Top-level output directory to populate as directed by
templatepath.`)
	templatePath string
)

var (
	templateText string
)

type (
	recursionContext struct {
		// XX  Copies of general args to command
		verbose bool
		format  string

		kvpArgs    []internal.KvpArg
		templLines int

		// Immutable after compilation of this file.
		//    https://golang.org/pkg/text/template/#hdr-Arguments
		tmpl *template.Template

		// Immutable after parsing of command line.
		// Selects which of 'kvpArgs' is exposed in name of output pathname
		//splitBaseFile,
		splitBaseDir []string // split at each '*inKeySeparator'

		// mutating state, changes at each iteration step
		//  X  Why a map rather than slice of K=V pairs?
		//       =>  Because template.Execute() doesn't understand
		//           the latter -- apparently reliant upon runtime type information.
		//         cf. https://golang.org/pkg/text/template/#Template.Execute
		//  X  Why indexed with base type 'string' rather than some
		//     defined type equivalent e.g. 'Key'?
		//       => Not acceptable to template.Execute():
		//              executing "singleton template" at <.CodeGenWarning>: can't
		//              evaluate field CodeGenWarning in type map[main.Key]string
		//         cf. https://golang.org/pkg/text/template/#hdr-Arguments

		// X  Provide template.Execute() with 'int' type if possible; otherwise 'string'.
		substitutions_var map[string]interface{}
	}
)

func UsageDump(helpAsMarkdown bool, cliUsage string) {
	toggleCode := func() { // XX  DRY
		if helpAsMarkdown {
			fmt.Fprintf(os.Stderr, "```\n")
		}
	}
	fmt.Fprintf(os.Stderr, "%s\n\n", usagePreamble)
	toggleCode()
	fmt.Fprintf(os.Stderr, "%s", cliUsage)
	fs.PrintDefaults()
	toggleCode()
}

func ParseArgs(genArgs []string, cliUsage string) (templLines int) {
	fs.Usage = func() {
		UsageDump(false, cliUsage)
		os.Exit(1)
	}

	usageWhy := func(why string) {
		UsageDump(false, cliUsage)
		fmt.Fprintf(os.Stderr, "gen command args: '%v'\n\n", genArgs)
		fmt.Fprintf(os.Stderr, "\n%s\n\n", why)
	}

	// X Flags required to precede all args other than the initial templatePath.
	//      After parsing, the arguments following the flags are available as
	//      the slice flag.Args() or individually as flag.Arg(i).
	//         https://golang.org/pkg/flag/#Args
	err := fs.Parse(genArgs)
	if err != nil {
		usageWhy(err.Error())
	}

	if nArgs := len(fs.Args()); nArgs < 1 {
		usageWhy("no path to template file found")
	} else if nArgs > 1 {
		usageWhy(fmt.Sprintf(
			"non-flag argument '%s' is not last arg on command line",
			fs.Args()[0]))
	}
	templatePath = fs.Args()[0]
	templLines, templateText = getTemplate(templatePath)
	return
}

func getTemplate(templatePath string) (int, string) {
	var err error
	var infile *os.File
	if infile, err = os.Open(templatePath); err != nil {
		log.Fatalf("Could not open input file \"%s\", err=%v",
			templatePath, err)
	}
	var stat os.FileInfo
	stat, err = infile.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	templateText := make([]byte, stat.Size())
	var nRead int
	if nRead, err = infile.Read(templateText); err != nil {
		log.Fatalf("Could not os.Read %s, err=%v", templatePath, err)
	}
	if nRead <= 0 {
		log.Fatalf("os.Read returned %d bytes", nRead)
	}
	templLines := internal.CountLines(infile)
	if err = infile.Close(); err != nil {
		log.Fatalf("Could not close %v, err=%v", infile, err)
	}
	return templLines, string(templateText)
}

func ExpandTemplate(verbose bool, format string, templLines int,
	kvpArgs []internal.KvpArg) {

	ctx := recursionContext{
		verbose:    verbose,
		format:     format,
		kvpArgs:    kvpArgs,
		templLines: templLines,
		splitBaseDir:      split(templatePath),
		substitutions_var: make(map[string]interface{}, 0),
	}

	if *inKeySeparator != "" {
		ctx.splitBaseDir = exciseChar(ctx.splitBaseDir)
	}

	var err error
	ctx.tmpl, err = template.New("" /*baseFile*/).Option("missingkey=error").
		Parse(templateText)
	if err != nil {
		log.Fatalln(err)
	}
	ctx.recurse(0)
}

func split(path string) []string {
	keysRE := regexp.MustCompile(`[a-zA-Z0-9_]+|[^a-zA-Z0-9_]+`)
	return keysRE.FindAllString(path, -1)
}

func exciseChar(frags []string) []string {
	//	separatorRE := regexp.MustCompile(`[` + *inKeySeparator + `]`)
	for ifrag := range frags {
		frags[ifrag] = strings.ReplaceAll(frags[ifrag], *inKeySeparator, "")
	}
	return frags
}

// X  Recursion here enumerates the combinations implied by the command-line
//    arguments K1=V11,V12,... K2=V21,V22,V23,... ...
//    This recursion is independent of any directory+file hierarchy specified
//    by 'templatePath'.
func (ctx *recursionContext) recurse(argIndex int) {
	// list of parameter values complete, so write out the file
	if argIndex == len(ctx.kvpArgs) {
		ctx.writeFile(argIndex)
		return
	}

	// Iterate from 'min' to 'max' for this 'argIndex' (and recursion level).
	eArg := &ctx.kvpArgs[argIndex]
	for _, enumVal := range eArg.Values {
		// XX  Document this data type conversion, and its effect on output.
		if intV, err := strconv.Atoi(enumVal); err == nil {
			ctx.substitutions_var[eArg.Key] = intV
		} else {
			ctx.substitutions_var[eArg.Key] = enumVal
		}

		ctx.recurse(argIndex + 1)
	}
}

func (ctx *recursionContext) substituteNames(splits []string) (
	fragmentsSubstituted []string, err error) {

	substitutions := 0
	for _, field := range splits {
		val, ok := ctx.substitutions_var[field]
		if !ok {
			fragmentsSubstituted = append(fragmentsSubstituted,
				field)
			continue
		}

		var vStr string
		switch v := val.(type) {
		case int:
			vStr = strconv.Itoa(v)
		case string:
			vStr = v
		}
		fragmentsSubstituted = append(fragmentsSubstituted,
			fmt.Sprintf(ctx.format, field, vStr))
		substitutions++
	}
	if substitutions == 0 && len(splits) > 0 {
		err = fmt.Errorf("WARNING: no substitutions made for pattern %v",
			splits)
	}
	return
}

func (ctx *recursionContext) writeFile(argIndex int) (outDir string) {

	buildOutPath := func(fragments []string) (segment string) {
		fragmentsSubstituted, err := ctx.substituteNames(fragments)
		if err != nil {
			log.Println(err)
		}
		segment = strings.Join(fragmentsSubstituted, "")
		return
	}

	outPathnameBottom := buildOutPath(ctx.splitBaseDir)
	outBottomDir := path.Dir(outPathnameBottom)
	outDir = *outTopDir + "/" + outBottomDir

	if err := os.MkdirAll(outDir, 0750); err != nil {
		log.Fatal(err)
	}

	// Make these synthetic K=V pairs available to the template.
	// XX  Which are useful?  How to document?
	ctx.substitutions_var["thisDir"] = path.Clean(outDir)
	//ctx.substitutions_var["parentDir"] = path.Dir(path.Clean(outDir))

	if ctx.verbose {
		log.Printf("Combination map:\n%s", ctx.formatMap())
	}

	outBaseName := path.Base(outPathnameBottom)
	outPath := outDir + "/" + outBaseName

	_, err := os.Stat(outPath)
	if err == nil {
		if !*clobber {
			log.Fatalf("Output file already exists: '%s'", outPath)
		}
		// X  If already existing, allow truncation by os.Create(), but no other
		//    operations.
		_ = os.Chmod(outPath, 0600)
	}
	outFile, err := os.Create(outPath)

	if err != nil {
		log.Fatalf("Failed to create output file '%s', err='%v'",
			outPath, err)
	}
	if err := ctx.tmpl.Execute(outFile, ctx.substitutions_var); err != nil {
		fmt.Fprintf(os.Stderr, "Template.Execute(outfile, map) returned  err=\n   %v", err)
		fmt.Fprintf(os.Stderr, "Contents of failing map:\n%s", ctx.formatMap())
		log.Fatalln("FATAL")
	}
	// X  Turn off 'w' bits, as a reminder to later readers of the output that file
	//    should not be edited.
	if err := outFile.Chmod(0440); err != nil {
		log.Fatal(err)
	}

	outLines := internal.CountLines(outFile)
	wrongCount := func() {
		log.Fatalf("outLines(%d) != ctx.templLines(%d), outPath=%s, templatePath=%s",
			outLines, ctx.templLines, outPath, templatePath)
	}

	// should line count disagreement be fatal error?
	if outLines > ctx.templLines {
		wrongCount()
	} else if outLines < ctx.templLines {
		wrongCount()
	}
	if err := outFile.Close(); err != nil {
		log.Fatal(err)
	}
	return
}

func (ctx *recursionContext) formatMap() (out string) {
	for k, v := range ctx.substitutions_var {
		out += fmt.Sprintf("   % 20s %15T '%v'\n", k, v, v)
	}
	return
}
