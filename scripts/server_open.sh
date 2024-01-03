#!/bin/bash

START=$SECONDS
dur=0
time_limit=500

source ../.env

host=$SERVER_HOST
port=$SERVER_PORTS

resp="$(curl -s http://localhost:8080)"

while [[ ${#resp} -eq 0 ]]
do
    resp="$(curl -s http://localhost:8080)"
    sleep 1
    dur=$((SECONDS - start))
    if [[ $dur -gt $time_limit ]]
    then
        echo "server connection failed by bash script after $dur seconds"
        exit
    fi
done

xdg-open http://$host:$port