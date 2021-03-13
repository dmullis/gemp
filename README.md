# gemp
=======
<!-- Copyright 2020 Donald Mullis. All rights reserved.
     https://github.github.com/gfm/
  -->


<!-- [Donald Mullis](https://github.com/dmullis)
 -->

# gemp: A Modest Template Processor

Gemp reads pairs specifying Key-to-list-of-Value mappings K=V1,V2...Vn,
and stores them for reformatting as directed by a specific command.

Gemp provides two commands, *gen* and *dump*,
each with various uses.

For trivial examples, read on.
## gen
### Push Key-Value Pair Combinations through Source Code Template

Expand template file according to all possible combinations of K=V1,V2... pairs,
write into multiple output files, naming each according to its
specific combination of values:

```sh
     $ echo '! Generated for {{.Color}} at {{.TimeHMS}}' >stamp+Color+.sh
     $ gemp -format '-%.0s%s' Color="Blue,Red" TimeHMS="$(date +%H:%M:%S)," \
           gen -inkeyseparator '+' -outtopdir . stamp+Color+.sh
     $ more stamp-*.sh
     ::::::::::::::
     stamp-Blue.sh
     ::::::::::::::
     ! Generated for Blue at 21:04:42
     ::::::::::::::
     stamp-Red.sh
     ::::::::::::::
     ! Generated for Red at 21:04:42
```

Format of the template file is as
specified in a Go standard library [template file](https://golang.org/pkg/text/template).
Usable functionality from ```template``` is constrained by gemp's limitation that the data structure
presented to the standard library's *Execute()* method is only a list of Key=Value
pairs, both Key and Value of string type.
A call is made to *Execute()* for each Value in the comma-separated value list Value+.
## dump
### Matching Constant Definitions Across Languages

Format all Key=Value+ pairs as specified on command line and write to standard output.
A pattern in the format of Go's [```fmt```](https://golang.org/pkg/fmt) package
take the place of the template file.

```sh
     $ KV=User=$USER\ TimeHMS="$(date +%H:%M:%S)"
     $ gemp -format 'export const %s = "%s";'$'' $KV dump
     export const User = "dmullis";
     export const TimeHMS = "20:41:06";
     $ gemp -format 'const %s = "%s"'$'' $KV dump
     const User = "dmullis"
     const TimeHMS = "20:41:06"
```
### Usage

[Generic](doc/usage.html) to both *gen* and *dump* commands.

[Specific to *gen*](doc/gen-usage.html).

If generating program source code, two difficulties may appear:
 1. For a satisfactory experience when debugging stack traces,
template expansions must match the number of lines in the template source code.
Workaround examples may be found in [```_test_src/```](./_test_src/).
 2. [gofmt](https://golang.org/cmd/gofmt/) is confused by template syntax e.g. "{{...}}".

"gemp" is a portmanteau of "Go-tEMPlate".

### See also

Other Go-based templating utilities, targeting somewhat different use cases:
 - [stringer](https://pkg.go.dev/golang.org/x/tools@v0.1.0/cmd/stringer)
 - [Kubernetes templates](https://pkg.go.dev/k8s.io/kubernetes/pkg/kubectl/util/templates)

