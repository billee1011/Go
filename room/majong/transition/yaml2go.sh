#!/bin/sh
#将当前目录下的yaml文件转换成go文件
gofileName="transition_cfg.go"
rm ${gofileName}
touch ${gofileName}
echo "package transition" >> ${gofileName}
SSID='`'
echo "var transitionCfg = ${SSID}" >> ${gofileName}
for file in ./*
do
    if test -f $file -a "${file##*.}"x = "yaml"x #遍历当前目录下所有yaml文件
    then
        IFS_old=$IFS
        while read line;
        do
            IFS=$'\n'
            echo "$line" >> ${gofileName}
        done < $file
        IFS=$IFS_old
    fi
done
echo "${SSID}" >> ${gofileName}
