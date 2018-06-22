go build -o bin/serviceloader steve/serviceloader 
go build -o bin/room/room.so -buildmode=plugin steve/room
go build -o bin/gateway/gateway.so -buildmode=plugin steve/gateway
go build -o bin/match/match.so -buildmode=plugin steve/match 
go build -o bin/login/login.so -buildmode=plugin steve/login 
