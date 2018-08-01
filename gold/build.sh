#!/bin/bash

# packet name
NAME=gold

mkdir ../bin/$NAME
go build  -o ./bin/$NAME.so -buildmode=plugin steve/$NAME
cp ./config.yml ../bin/$NAME/
cp ./config.yml  ./bin/$NAME/
cp ./bin/$NAME.so ../bin/$NAME/



