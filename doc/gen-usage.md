command 'gen' usage:

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


```
gemp [-format=%-.s-%s] [-h=false] [-helpAsMarkdown=false] [-kvpluspath=] [-verbose=false] [K=V1,V2...Vn]* (gen [flags] input_file | dump)

  -clobber
    	Overwrite already-existing output files.
  -inkeyseparator string
    	Input files may be visually distinguished from output
    	files they generate by inclusion of a specified character.  The character
    	must not be legal in a Go identifer ([a-zA-Z0-9_]).  Any instances of
    	the character will be omitted from output file names.
    	
    	Candidates for '-inkeyseparator' usage must seek a compromise:
    	   a. Escape special treatment by build tools, command shells,
    	      or GNU's 'readline' library, and
    	   b. Not collide with other non-alphanums wanted within filenames.
    	A few non-alphanumeric candidates: + ~ @  %
  -outtopdir string
    	Top-level output directory to populate as directed by
    	templatepath. (default ".")
```
