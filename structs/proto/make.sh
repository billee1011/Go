
ROOT=../../proto 
DEST=.

function makeProto() {
    local dirName=$1
    local relativeDirName=${dirName#$ROOT/}

    if [ -z "`ls $dirName | grep .*\.proto`" ]; then 
        echo $dirName "不包含proto文件"
    else 
        echo 正在处理: $dirName 

        # $ROOT/.. 表示 ROOT 目录的父目录， 默认将该目录加入 import 搜索列表
        protoc -I $dirName -I $ROOT/.. --go_out=plugins=grpc:$DEST/$relativeDirName $dirName/*.proto

        # 替换 import xxx "proto/yyy" 为  import xxx "steve/structs/proto/yyy"
        sed -i 's/"proto\//"steve\/structs\/proto\//g' $DEST/$relativeDirName/*.go 
    fi 

    for file in $dirName/*
    do 
        if test -d $file 
        then 
            local subDir=${file#$ROOT/}
            if !(test -d $subDir)
            then 
                mkdir $DEST/$subDir
            fi 
            makeProto $file 
        fi 
    done 
}

makeProto $ROOT