#! /bin/sh

set -o nounset
set -e

usage () {
    set +x
    printf "\n%s\n" "$1"
    exit 1
}

set -v
# X  Sanitizing of command-line args -- Go-specific security measure?
#    Note that environment variables are similarly "sanitized" or
#    more literally perhaps "quoted".
#
# X  Command line '\t' is turned into '\\t' as appearing in os.Args[] -- before
#    execution flag.Parse().
#    Observed that /proc/PID/cmdline holds '\t'.
#
# X  Bash command line $'\t' turned into '\t' as appearing in os.Args[]
#    /proc/PID/cmdline holds '	' i.e. the TAB character.
#
# X  Sh command line '	' i.e. literal TAB character is turned into '\t' as
#    appearing in os.Args[] -- the desired effect here.

echo Go or JavaScript constant definitions
gemp -format '	const %s = "%s"' 'K=V' dump
gemp -format '	const %s = %s' 'K=1' dump
gemp -format '	const %s = %s' 'K=1' dump

echo Equivalent invocations:
gemp -format '%s:%s' 'K=V1,V2' dump
echo 'K=V1,V2' | gemp -format '%s:%s' -kvpluspath /dev/stdin dump
