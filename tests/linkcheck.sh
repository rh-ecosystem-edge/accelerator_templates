#!/bin/bash


## use markdownlint from https://github.com/igorshubovych/markdownlint-cli

ERR=0
PODMAN=podman
PATTERN=$@

if [ -z "$PATTERN" ]; then
    PATTERN=$(find . -path "./src" -prune -o  -name "*.md" -print )
fi

for arg in $PATTERN
do
    #echo $arg
    CHECKFILES=""

    if [[ -f "${arg}" &&  "${arg}" =~ .md$ ]]; then
        CHECKFILES=${arg}
    fi
    if [[ -d "${arg}" ]]; then
        CHECKFILES=$(find ${arg} -name "*.md" -print)
    fi
    if [ -z "${CHECKFILES}" ]; then
        continue
    fi
    #echo "MATCH!"

    $PODMAN run -v ${PWD}:/tmp:ro  --workdir /tmp --rm -i ghcr.io/tcort/markdown-link-check:stable --quiet $CHECKFILES
    E=$?
    if [ $E -ne 0 ]; then
        ERR=$E
    fi
done
exit $ERR
