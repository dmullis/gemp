Usage:
```
gemp [-format=%-.s-%s] [-h=false] [-helpAsMarkdown=false] [-kvpluspath=] [-verbose=false] [K=V1,V2...Vn]* (gen [flags] input_file | dump)

  -format string
    	Format string syntax is that of Go's 'fmt' package, with exactly
    	two string expansion codes e.g. "%s-%s" required.
    	
    	Each pair of Key, Value strings is passed to fmt.Printf, along with this
    	format string.
    	  'gen'  Result is reinserted into each file's output pathname
    	  'dump' Results written line-by-line to stdout.
    	Note that prefixing with '%-.s', drops a string from output.
    	 (default "%-.s-%s")
  -h	Repeat this message.
  -helpAsMarkdown
    	Format help output, if any, as Markdown
  -kvpluspath string
    	Alternative to specifying K=V+ pairs on the command line. Arg is a
    	path to an input file containing Key=Value+ pairs, in 'sh' syntax.
    	Lines of commentary, beginning with '#', are ignored.
  -verbose
    	Log heavily
  (K=V1,V2...Vn)*
        Any number of Key=Value+ pairs, where Value+ may be a comma-
	separated list of multiple string values to be substituted serially
	into each of multiple output directories or files.
```

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
