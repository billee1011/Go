#!/bin/bash

i=1;
i1=10000000;
i2=10000000;
i3=10000000;
i4=10000;
echo "2222"
while [ $i -le 50 ]
do
   echo $i
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player (playerID,showUID,accountID,type,nickname) values (${i1},${i2},${i3},2,'robot${i1}');"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_currency (playerID,coins) values (${i1},${i4});"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_game (playerID,gameID,winningRate) values (${i1},1,50);"
    d=$(date +%M-%d\ %H\:%m\:%S)
    echo "INSERT $i @@ $d"    
    i=$(($i+1))
    i1=$(($i1+1))
    i2=$(($i2+1))
    i3=$(($i3+1))	

    sleep 0.05
done
i=1;
while [ $i -le 50 ]
do
   echo $i
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player (playerID,showUID,accountID,type,nickname) values (${i1},${i2},${i3},2,'robot${i1}');"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_currency (playerID,coins) values (${i1},${i4});"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_game (playerID,gameID,winningRate) values (${i1},2,50);"
    d=$(date +%M-%d\ %H\:%m\:%S)
    echo "INSERT $i @@ $d"    
    i=$(($i+1))
    i1=$(($i1+1))
    i2=$(($i2+1))
    i3=$(($i3+1))	

    sleep 0.05
done
i=1;
while [ $i -le 50 ]
do
   echo $i
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player (playerID,showUID,accountID,type,nickname) values (${i1},${i2},${i3},2,'robot${i1}');"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_currency (playerID,coins) values (${i1},${i4});"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_game (playerID,gameID,winningRate) values (${i1},3,50);"
    d=$(date +%M-%d\ %H\:%m\:%S)
    echo "INSERT $i @@ $d"    
    i=$(($i+1))
    i1=$(($i1+1))
    i2=$(($i2+1))
    i3=$(($i3+1))	

    sleep 0.05
done
i=1;
while [ $i -le 50 ]
do
   echo $i
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player (playerID,showUID,accountID,type,nickname) values (${i1},${i2},${i3},2,'robot${i1}');"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_currency (playerID,coins) values (${i1},${i4});"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_game (playerID,gameID,winningRate) values (${i1},4,50);"
    d=$(date +%M-%d\ %H\:%m\:%S)
    echo "INSERT $i @@ $d"    
    i=$(($i+1))
    i1=$(($i1+1))
    i2=$(($i2+1))
    i3=$(($i3+1))	

    sleep 0.05
done
i=1;
while [ $i -le 50 ]
do
   echo $i
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player (playerID,showUID,accountID,type,nickname) values (${i1},${i2},${i3},2,'robot${i1}');"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_currency (playerID,coins) values (${i1},10000);"
    mysql -h192.168.9.200 -uroot -p'SDF-XSP-0056' 192_168_9_162_player -e "insert into t_player_game (playerID,gameID,winningRate) values (${i1},4,50);"
    d=$(date +%M-%d\ %H\:%m\:%S)
    echo "INSERT $i @@ $d"    
    i=$(($i+1))
    i1=$(($i1+1))
    i2=$(($i2+1))
    i3=$(($i3+1))	

    sleep 0.05
done
exit 0
