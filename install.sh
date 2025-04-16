#!/bin/sh

set -e

rm -rf build
mkdir build

go build -o build/dos2unix ./cmds/dos2unix && cp build/dos2unix ~/bin
go build -o build/unix2dos ./cmds/unix2dos && cp build/unix2dos ~/bin
