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
