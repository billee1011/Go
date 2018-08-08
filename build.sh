#!/bin/bash

go build -o bin/serviceloader steve/serviceloader 
go install steve/serviceloader

go build -o bin/room/room.so -o room/room.so -buildmode=plugin steve/room
cp room/config.yml configs/room/config.yml 

go build -o bin/gateway/gateway.so -o gateway/gateway.so -buildmode=plugin steve/gateway
cp gateway/config.yml configs/gateway/config.yml 


go build -o bin/match/match.so -o match/match.so -buildmode=plugin steve/match 
cp match/config.yml configs/match/config.yml 

go build -o bin/login/login.so -o login/login.so -buildmode=plugin steve/login 
cp login/config.yml configs/login/config.yml 


go build -o bin/hall/hall.so -o hall/hall.so -buildmode=plugin steve/hall 
cp hall/config.yml configs/hall/config.yml 

go build -o bin/gold/gold.so -o gold/gold.so -buildmode=plugin steve/gold 
cp gold/config.yml configs/gold/config.yml 

go build -o bin/robot/robot.so -o robot/robot.so -buildmode=plugin steve/robot 
cp robot/config.yml configs/robot/config.yml 

go build -o bin/configuration/configuration.so -o configuration/configuration.so -buildmode=plugin steve/configuration 
cp configuration/config.yml configs/configuration/config.yml 

go build -o bin/configuration/configuration.so -o configuration/configuration.so -buildmode=plugin steve/configuration 
cp configuration/config.yml configs/configuration/config.yml 

go build -o bin/msgserver/msgserver.so -o msgserver/msgserver.so -buildmode=plugin steve/msgserver 
cp msgserver/config.yml configs/msgserver/config.yml

go build -o bin/alms/alms.so -o alms/alms.so -buildmode=plugin steve/alms 
cp alms/config.yml configs/alms/config.yml

./simulate/packtests.sh 


if [ "$1"="pack" ];then  
    tar -czf server.tar.gz bin configs
fi

