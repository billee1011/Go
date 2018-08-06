pushd gateway 
nohup serviceloader gateway --config=config.yml  &
popd 

pushd room 
nohup serviceloader room --config=config.yml  &
popd 


pushd hall 
nohup serviceloader hall --config=config.yml  &
popd 

pushd login 
nohup  serviceloader login --config=config.yml  &
popd 

pushd match 
nohup  serviceloader match --config=config.yml  &
popd 

pushd gold 
sh ./start.sh
popd 

pushd robot 
nohup serviceloader robot --config=config.yml  &
popd

pushd gold
nohup serviceloader gold --config=config.yml  &
popd


pushd msgserver
nohup serviceloader msgserver --config=config.yml  &
popd


