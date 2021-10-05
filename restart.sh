#!/bin/bash

path=$PWD

git pull
killall app

cd $path/server
nohup go run app.go &

sleep 3
num=$(netstat -nlpt | grep ":::23333" | grep -v "grep" | wc -l)
if [ $num -eq 1 ]; then
    echo "run app success"
else
    cd $path/server
    tail -n 10 nohup.out
    echo "run app fail"
fi