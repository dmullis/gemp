#! /bin/sh

# X  Why does this not remove executable in $GOBIN/ ?
#      https://golang.org/cmd/go/#hdr-Remove_object_files_and_cached_files
go clean -i

rm -f $(go env GOBIN)/gemp

