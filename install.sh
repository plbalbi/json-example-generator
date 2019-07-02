#!/bin/bash

function logFailureIfNecessary {
    # $1 Message
    if [[ $? -ne 0 ]]; then
        echo "FAILURE: $1"
        exit 1
    fi
}

go get golang.org/x/tools/cmd/goyacc
go get github.com/golang-collections/collections/queue
logFailureIfNecessary "Failed to get go dependencies"

cd parser

goyacc -l -o parser.go parser.y
logFailureIfNecessary "Failed to compile parser from 'parser.y' file"
