# protoc -I=${GOPATH}/src:. --gogofaster_out=${GOPATH}/src  steve/clientpb/game.proto
# protoc -I=${GOPATH}/src:. --gogofaster_out=${GOPATH}/src  steve/clientpb/msgid/msgid.proto

protoc -I ./msgid --go_out=./msgid ./msgid/*.proto  
protoc -I ./ --go_out=./ ./*.proto  