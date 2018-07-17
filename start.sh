pushd gateway 
serviceloader gateway --config=config.yml  &
popd 

pushd room 
serviceloader room --config=config.yml  &
popd 


pushd hall 
serviceloader hall --config=config.yml  &
popd 

pushd login 
serviceloader login --config=config.yml  &
popd 

pushd match 
serviceloader match --config=config.yml  &
popd 
