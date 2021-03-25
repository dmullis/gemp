<!-- Copyright 2020 Donald Mullis. All rights reserved.
     https://github.github.com/gfm/
  -->


<!-- [Donald Mullis](https://github.com/dmullis)
 -->

# gemp: A Recursive CLI Expander of Go Template Files

Gemp reads pairs specifying Key-to-list-of-Value mappings K=V1,V2...Vn,
and stores them for reformatting as directed by a specific command.

Use cases are generation of source code variants differing only by simple substitution,
but possibly many combinations of values.
Substitutions are performed by Go's ```text/template``` standard library.

For trivial examples, read on.
