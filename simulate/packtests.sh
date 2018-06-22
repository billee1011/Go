# 打包测试二进制文件，提供给测试人员

export OUT_DIR=bin/

go test steve/simulate/tests/logintests -run Test_Login -c -o $OUT_DIR/login.test
go test steve/simulate/tests/scxltests -run Test_StartGame -c -o $OUT_DIR/startgame.test
