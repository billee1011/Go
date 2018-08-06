# 查找替换
find . -name '*.yml' | xargs perl -pi -e 's|127.0.0.1|192.168.7.108|g'
find . -name '*.yml' | xargs perl -pi -e 's|192.168.8.17|192.168.7.108|g'
find . -name '*.yml' | xargs perl -pi -e 's|root|backuser|g'
find . -name '*.yml' | xargs perl -pi -e 's|123456|Sdf123esdf|g'


