#!/bin/sh
#将当前目录下的yaml文件转换成go文件
for file in ./*
do
    if test -f $file -a "${file##*.}"x = "yaml"x #遍历当前目录下所有yaml文件
    then
        basefileName=`basename $file  .yaml` 
        gofileName=${basefileName}"_cfg.go"
        rm ${gofileName}
        touch ${gofileName}
        echo "package ${basefileName}" >> ${gofileName}
        SSID='`'
        echo "var ${basefileName}Cfg = ${SSID}" >> ${gofileName}
        IFS_old=$IFS
        while read line;
        do
            IFS=$'\n'
            echo "$line" >> ${gofileName}
        done < $file
        IFS=$IFS_old
        echo "${SSID}" >> ${gofileName}
    fi
done