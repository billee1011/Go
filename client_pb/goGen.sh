protoc -I ./msgId --go_out=./msgId ./msgId/*.proto
protoc -I ./room --go_out=./room ./room/*.proto
protoc -I ./login --go_out=./login ./login/*.proto 
protoc -I ./gate --go_out=./gate ./gate/*.proto 
