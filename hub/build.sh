#!/bin/bash

program="hub_server"

export GOPATH=`pwd`

if [ -d output ]; then
rm -rf output
fi

echo "begin build..."
go build -o ./$program src/main.go

mkdir output
mkdir -p output/bin
mkdir -p output/conf

mv ./${program} ./output/bin/
cp -r conf/ output/

echo "build success."
