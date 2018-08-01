#!/bin/bash

# packet name
NAME=room

mkdir ../release/$NAME
echo "begin building..."
go build  -o ./$NAME.so -buildmode=plugin steve/$NAME

echo "begin copy yml"
cp  -f ./config.yml ../release/$NAME/

echo "begin copy start.sh, stop.sh"
cp -f ./start.sh ../release/$NAME/
cp -f ./stop.sh ../release/$NAME/

echo "begin cp so"
cp -f  ./$NAME.so ../release/$NAME/

echo "end  build..."



