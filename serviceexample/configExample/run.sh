consul kv put product/version 1.0.0
consul kv put product/1.0.0/game/name majong
consul kv put product/1.0.0/game/version 2.0.1
consul kv put product/1.0.0/hello hello 

go build -buildmode plugin 
serviceloader configExample
rm *.so
