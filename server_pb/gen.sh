protoc -I majong --go_out=plugins=grpc:majong majong/*.proto  
protoc -I user --go_out=plugins=grpc:user user/*.proto
protoc -I gateway --go_out=plugins=grpc:gateway gateway/*.proto  
protoc -I room --go_out=plugins=grpc:room room/*.proto  
protoc -I match --go_out=plugins=grpc:match match/*.proto  