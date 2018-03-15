
go build -buildmode plugin 
serviceloader loggerexample --config=./config.yml
rm *.so
