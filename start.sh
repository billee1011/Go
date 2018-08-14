#!/bin/bash


function startserver() {
    path=$1
    name=$2
    cname=$3
    pushd $path 
    serviceloader $name --config=config.yml &
    # x=`consul catalog services | grep $cname | wc -l`
    # while [[ $x -eq 0 ]]; do
    #     echo 等待 $name 启动完成$x

    #     sleep 1
    #     x=`consul catalog services | grep $cname | wc -l` 
    # done 
    # sleep 1
    echo $name 启动完成$x
    popd 
}
startserver configuration configuration configuration
startserver gateway gateway gate
startserver room room room
startserver hall hall hall
startserver login login login
startserver robot robot robot
startserver gold gold gold
startserver msgserver msgserver msgserver
startserver alms alms alms
startserver match match match

pushd back
serviceloader back --config=config.yml &
popd