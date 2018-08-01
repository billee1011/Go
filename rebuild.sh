#!/bin/bash

BIN=Release


pushd serviceloader
./build.sh
popd


pushd room
./build.sh
popd


pushd gateway
./build.sh
popd

pushd match
./build.sh
popd

pushd login
./build.sh
popd

pushd hall
./build.sh
popd

pushd robot
./build.sh
popd


./simulate/packtests.sh 


if [ "$1"="pack" ];then  
    tar -czf server.tar.gz release
fi

