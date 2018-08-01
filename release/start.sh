
echo "gateway------------------------------------\n"
pushd gateway 
nohup serviceloader gateway --config=config.yml  &
popd 

echo "room------------------------------------\n"
pushd room 
nohup serviceloader room --config=config.yml  &
popd 

echo "hall----------------------------------\n"
pushd hall 
nohup serviceloader hall --config=config.yml  &
popd 

echo "login---------------------------------\n"
pushd login 
nohup serviceloader login --config=config.yml  &
popd 

echo "match---------------------------------\n"
pushd match 
nohup serviceloader match --config=config.yml  &
popd 
echo "gold---------------------------------\n"
pushd gold 
sh ./start.sh
popd

./p.sh

