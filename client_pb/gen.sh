PROTO_PATH=~/steve.protos/client_pb 

mkdir msgId -p 
protoc -I $PROTO_PATH/msgId --go_out=./msgId $PROTO_PATH/msgId/*.proto 

mkdir room -p 
protoc -I $PROTO_PATH/room --go_out=./room $PROTO_PATH/room/*.proto 

mkdir login -p 
protoc -I $PROTO_PATH/login --go_out=./login $PROTO_PATH/login/*.proto 

mkdir gate -p
protoc -I $PROTO_PATH/gate --go_out=./gate $PROTO_PATH/gate/*.proto 
