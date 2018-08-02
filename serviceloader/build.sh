#!/bin/bash


mkdir ../release

go build -o ../release/serviceloader steve/serviceloader
go install steve/serviceloader

echo "end  build..."



