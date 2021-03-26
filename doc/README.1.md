<!-- Copyright 2020 Donald Mullis. All rights reserved.
     https://github.github.com/gfm/
  -->


<!-- [Donald Mullis](https://github.com/dmullis)
 -->

# gemp: A Recursive CLI Expander of Go Template Files

Gemp reads pairs specifying Key-to-list-of-Value mappings ```K=V1,V2...Vn```,
and stores them for use according to the purpose of a specific subcommand
```gen``` or ```dump```.

Typical use cases are generation of source code variants differing only by simple substitution,
but possibly many combinations of values.
Substitutions are performed by Go's ```text/template``` or ```fmt``` standard library packages.

For trivial examples, read on.
