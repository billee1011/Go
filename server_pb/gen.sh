# protoc -I majong --go_out=plugins=grpc:majong majong/*.proto  
protoc -I login --go_out=plugins=grpc:login login/*.proto
protoc -I user --go_out=plugins=grpc:user user/*.proto
protoc -I gateway --go_out=plugins=grpc:gateway gateway/*.proto  
protoc -I room_mgr --go_out=plugins=grpc:room_mgr room_mgr/*.proto  
protoc -I match --go_out=plugins=grpc:match match/*.proto
# protoc -I ddz --go_out=plugins=grpc:ddz ddz/*.proto
protoc -I robot --go_out=plugins=grpc:robot robot/*.proto
