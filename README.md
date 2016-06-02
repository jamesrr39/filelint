Duplicate file finder. Matches on length, then on md5 hash.

[![Build Status](https://travis-ci.org/jamesrr39/filelint.svg?branch=master)](https://travis-ci.org/jamesrr39/filelint)

Build

    go build -o filelint

Run

    ./filelint [-v] DIRECTORY

example:

    ./filelint ~/Documents

Test

    go test ./...
