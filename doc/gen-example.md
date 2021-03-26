### Push Key-Value Pair Combinations through Source Code Template

Expand some template file according to all possible combinations of K=V1,V2... pairs,
write into multiple output files, naming each according to its
specific combination of values.
Format of the template file is as
specified in a Go standard library [template file](https://golang.org/pkg/text/template).
Usable functionality from
```template```
is constrained by gemp's limitation that the data structure
presented to the
```template.Execute()```
is simply a list of Key=Value
pairs, both Key and Value of string type.

```
     # Write out a file containing valid ```template``` syntax
     $ echo '! Generated for {{.Color}} at {{.TimeHMS}}' >stamp+Color+.sh
     # Expand into two output files containing substitutions.
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

In the example above, a call is made to ```Execute()``` for each
of ```Blue``` and ```Red``` from the comma-separated list ```Color```.

