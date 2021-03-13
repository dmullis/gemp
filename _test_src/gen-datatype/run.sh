#! /bin/sh

set -ue

usage () {
    set +x
    printf "%s\n" "$1"
    printf "Examples:\n"
    printf "\t%s %s\n" $0 \
           "[ byte | float64 | 'interface{f()}' | 'map[string]int' | 'struct{}' | 'struct{a int; b [9]bool}' | ... ]"
    exit 1
}

if [ $# != 1 ]
then
    usage "Exactly one arg expected, a Go type definition."
fi
goDataType="$1"

EXAMPLE_BASENAME=type-parameter
TEMPLATENAME=${EXAMPLE_BASENAME}+GoDataType+.go
OUTFILE_SEPARATOR="__"

(
    set -x
    gemp -format="${OUTFILE_SEPARATOR}%.0s%s" \
       Copyright="Copyright 2020 ${USER}. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file." \
       CodeGenWarning="Generated code -- source was ${TEMPLATENAME}." \
       NewLine="
" \
       GoDataType="${goDataType}" \
       gen \
       -inkeyseparator '+' \
       -clobber \
       -outtopdir . \
       ./${TEMPLATENAME}
)

(
    printf "\nNumber of lines should be same in both source -- \"%s\" -- and
generated files, for ease of debugging with symbolic backtraces:\n\n" \
           ${TEMPLATENAME}
    set -x
    wc --lines *.go | grep -v total
)

printf "\nTo run generated code:\n"
for f in ${EXAMPLE_BASENAME}${OUTFILE_SEPARATOR}*.go
do
    printf "\t%s\n" "go run \"./${f}\""
done
