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
presented to the standard library's ```Execute()``` method is only a list of Key=Value
pairs, both Key and Value of string type.
A call is made to ```Execute()``` for each Value in the comma-separated value list Value+.
