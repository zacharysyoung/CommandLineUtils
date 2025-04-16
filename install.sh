#!/bin/sh

set -e

rm -rf build
mkdir build

for cmd in \
    dos2unix lspath searchup tree unix2dos
do
    echo $cmd
    go build -o build/$cmd ./cmds/$cmd && cp build/$cmd ~/bin
done
