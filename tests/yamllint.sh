#!/bin/bash

## run aspell over all markdown files listed on the command line
## if the argument is a directory aspell all the markdown file in that directory tree

ERR=0
LINT_CONFIG=./tests/yamllint.cfg
LINTER=/usr/bin/yamllint
PATTERN=$@

if [ -z "$PATTERN" ]; then
    PATTERN=$(find . -path "./src" -prune -o  -name "*.y[a]ml"  -print) 
fi

for arg in $PATTERN
do
    #echo $arg
    CHECKFILES=""

    if [[ -f "${arg}" && "${arg}" =~ .y[a]ml$ ]]; then
        CHECKFILES=${arg}
    fi
    if [[ -d "${arg}" ]]; then
        CHECKFILES=${arg}
    fi
    if [ -z "${CHECKFILES}" ]; then
        continue
    fi
    #echo "MATCH!"

    for file in ${CHECKFILES}
    do
        $LINTER -c ${LINT_CONFIG} ${file}
        if [ $? -ne 0 ]; then
            ERR=1
        fi
    done
done
exit $ERR
