
ROOT=../../proto 
DEST=.

function makeProto() {
    local dirName=$1
    local relativeDirName=${dirName#$ROOT/}

    echo "current path:" $relativeDirName 
    protoc -I $dirName --go_out=plugins=grpc:$DEST/$relativeDirName $dirName/*.proto

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

