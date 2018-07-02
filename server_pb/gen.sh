protoc -I majong --go_out=plugins=grpc:majong majong/*.proto  
protoc -I user --go_out=plugins=grpc:user user/*.proto
protoc -I ddz --go_out=plugins=grpc:ddz ddz/*.proto