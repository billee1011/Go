protoc -I majong --gofast_out=plugins=grpc:majong majong/*.proto  
protoc -I user --gofast_out=plugins=grpc:user user/*.proto
protoc -I gateway --gofast_out=plugins=grpc:gateway gateway/*.proto  
protoc -I room --gofast_out=plugins=grpc:room room/*.proto  