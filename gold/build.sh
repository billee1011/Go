#!/bin/bash

# packet name
NAME=gold

mkdir ../bin/$NAME
echo "begin building..."
go build  -o ./bin/$NAME.so -buildmode=plugin steve/$NAME

echo "begin copy yml"
cp  -f ./config.yml ../bin/$NAME/
cp -f  ./config.yml  ./bin/

echo "begin copy start.sh, stop.sh"
cp -f ./start.sh ../bin/$NAME/
cp -f  ./start.sh  ./bin/
cp -f ./stop.sh ../bin/$NAME/
cp -f  ./stop.sh  ./bin/

echo "begin cp so"
cp -f  ./bin/$NAME.so ../bin/$NAME/

echo "end  build..."



