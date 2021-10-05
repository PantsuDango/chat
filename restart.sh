#!/bin/bash

git pull
killall chat

nohup go run chat.go &
sleep 3

num=$(netstat -nlpt | grep ":::23333" | grep -v "grep" | wc -l)
if [ $num -eq 1 ]; then
    echo "run chat success"
else
    tail -n 10 nohup.out
    echo "run chat fail"
fi