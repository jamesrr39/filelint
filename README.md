Duplicate file finder. Matches on length, then on md5 hash.

Build

    go build -o filelint

Run

    ./filelint [-v] DIRECTORY

example:

    ./filelint ~/Documents

Test

    go test ./...
