#!/bin/bash

echo "wget rscli"
cd /usr/local/bin/
wget http://192.168.8.189:1313/%E6%96%B0%E6%89%8B%E6%8C%87%E5%8D%97/rsctl%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/files/rscli 
chmod 755 /usr/local/bin/rscli

echo "make directory"
mkdir -p ~/.rscli

cd ~/.rscli

echo "wget rscli_ca.crt"
wget http://192.168.8.189:1313/%E6%96%B0%E6%89%8B%E6%8C%87%E5%8D%97/rsctl%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/files/rscli_ca.crt

echo "wget rscli.yml"
wget http://192.168.8.189:1313/%E6%96%B0%E6%89%8B%E6%8C%87%E5%8D%97/rsctl%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/files/rscli.yml 

cd -

echo "append 192.168.7.12 rsctl.fz.stevegame.red to /etc/hosts "
echo 192.168.7.12 rsctl.fz.stevegame.red >> /etc/hosts
