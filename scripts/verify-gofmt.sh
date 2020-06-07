#!/usr/bin/env bash

set -e

files=$(gofmt -l .)

if [[ -nz "$files" ]]; then
    echo
    echo 'The following files are unformatted:'
    echo "$files"
    echo
    echo 'Please run `go fmt ./...` to fix.'
    echo
    exit 1
fi

exit 0
