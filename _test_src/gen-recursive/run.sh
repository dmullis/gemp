#! /bin/sh

. ./lib.sh

# Execute the generated code.
for d in $(ls -t1 $TOPOUTDIR | tac)
do
    if ! [ -d $TOPOUTDIR/$d ]
    then
        continue
    fi
    echo
    (
        cd $TOPOUTDIR/$d
        go test -v -bench .
    )
done
