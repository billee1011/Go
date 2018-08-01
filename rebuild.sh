#!/bin/bash

BIN=Release


pushd serviceloader
./build.sh
popd


pushd room
./build.sh
popd


pushd gateway
sh ./build.sh
popd

pushd match
sh ./build.sh
popd

pushd login
sh ./build.sh
popd

pushd hall
sh ./build.sh
popd

pushd robot
sh ./build.sh
popd


sh ./simulate/packtests.sh


if [ "$1"="pack" ];then  
    tar -czf server.tar.gz release
fi

