#!/bin/bash

pushd configuration 
serviceloader configuration --config=config.yml &
popd 

# 其他服务启动依赖配置服
sleep 2

pushd gateway 
serviceloader gateway --config=config.yml  &
popd 

pushd room 
serviceloader room --config=config.yml  &
popd 


pushd hall 
serviceloader hall --config=config.yml  &
popd 

pushd login 
serviceloader login --config=config.yml  &
popd 

pushd match 
serviceloader match --config=config.yml  &
popd 

pushd robot 
serviceloader robot --config=config.yml  &
popd

pushd gold
serviceloader gold --config=config.yml  &
popd


pushd msgserver
serviceloader msgserver --config=config.yml  &
popd


