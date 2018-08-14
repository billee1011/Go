
echo "configuration------------------------------------\n"
pushd configuration
sh ./start.sh
popd

sleep 1s

echo "gold---------------------------------\n"
pushd gold
sh ./start.sh
popd

sleep 1s

echo "propserver---------------------------------\n"
pushd propserver
sh ./start.sh
popd

sleep 1s

echo "gateway------------------------------------\n"
pushd gateway
sh ./start.sh
popd 

sleep 1s

echo "room------------------------------------\n"
pushd room
sh ./start.sh
popd

sleep 1s

echo "hall----------------------------------\n"
pushd hall
sh ./start.sh
popd

sleep 1s

echo "login---------------------------------\n"
pushd login
sh ./start.sh
popd

sleep 1s

echo "match---------------------------------\n"
pushd match
sh ./start.sh
popd

sleep 1s

echo "robot---------------------------------\n"
pushd robot
sh ./start.sh
popd

sleep 1s

echo "msgserver---------------------------------\n"
pushd msgserver
sh ./start.sh
popd

sleep 1s

echo "mailserver---------------------------------\n"
pushd mailserver
sh ./start.sh
popd




sleep 1s

./p.sh

