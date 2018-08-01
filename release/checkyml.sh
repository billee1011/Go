// 查找替换
find . -name '*.yml' | xargs perl -pi -e 's|127.0.0.1|192.168.7.108|g'
