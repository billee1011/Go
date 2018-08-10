
echo "configuration------------------------------------\n"
pushd configuration
sh ./start.sh
popd

echo "gold---------------------------------\n"
pushd gold
sh ./start.sh
popd

echo "gateway------------------------------------\n"
pushd gateway 
#nohup serviceloader gateway --config=config.yml  &
sh ./start.sh
popd 

echo "room------------------------------------\n"
pushd room 
#nohup serviceloader room --config=config.yml  &
sh ./start.sh
popd 

echo "hall----------------------------------\n"
pushd hall 
#nohup serviceloader hall --config=config.yml  &
sh ./start.sh
popd 

echo "login---------------------------------\n"
pushd login 
#nohup serviceloader login --config=config.yml  &
sh ./start.sh
popd 

echo "match---------------------------------\n"
pushd match 
#nohup serviceloader match --config=config.yml  &
sh ./start.sh
popd 

echo "robot---------------------------------\n"
pushd robot
sh ./start.sh
popd

echo "msgserver---------------------------------\n"
pushd msgserver
sh ./start.sh
popd

./p.sh

