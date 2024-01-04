#!/bin/bash

START=$SECONDS
dur=0
time_limit=10

host="localhost"
port="8082"

resp="$(curl -s http://$host:$port)"

while [[ ${#resp} -eq 0 ]]
do
    resp="$(curl -s http://$host:$port)"
    sleep 1
    dur=$((SECONDS - start))
    if [[ $dur -gt $time_limit ]]
    then
        echo "server connection failed by bash script after $dur seconds"
        exit
    fi
done

go run cmd/sender/main.go