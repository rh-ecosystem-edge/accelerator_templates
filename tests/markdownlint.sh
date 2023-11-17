#!/bin/bash


## use markdownlint from https://github.com/igorshubovych/markdownlint-cli

ERR=0
PODMAN=podman
PATTERN=$@

if [ -z "$PATTERN" ]; then
    PATTERN=$(find . -path "./src" -prune -o  -name "*.md"  -print)
fi

for arg in $PATTERN
do
    #echo $arg
    CHECKFILES=""

    if [[ -f "${arg}" &&  "${arg}" =~ .md$ ]]; then
        CHECKFILES=${arg}
    fi
    if [[ -d "${arg}" ]]; then
        CHECKFILES=$(find ${arg} -name "*.md")
    fi
    if [ -z "${CHECKFILES}" ]; then
        continue
    fi
    #echo "MATCH!"

    for file in ${CHECKFILES}
    do

        $PODMAN  run -v $PWD:/workdir:z ghcr.io/igorshubovych/markdownlint-cli:latest --disable MD013 --ignore '*.!(md)'  --  $file
        if [ $? -ne 0 ]; then
            ERR=$?
        fi
    done
done
exit $ERR
