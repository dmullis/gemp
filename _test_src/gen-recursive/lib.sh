#! /bin/sh

set -ue

usage () {
    set +x
    printf "\n%s\n" "$1"
    exit 1
}

: ${TOPOUTDIR:=./test-recursive}   # X  "test/" is a standard name known to `go clean`

set -x
