#!/bin/bash

source ../.env

pport=$POSTGRES_PORTS
nport=$NATS_PORTS
sport=$SERVER_PORTS

echo "The script will clear the ports used by the program: $pport, $nport, $sport."
echo "This can be useful if the server is not running in a container."
echo "Allow? (y/n)"

read symbol

if [ $symbol = 'y' ] || [ $symbol = 'Y' ] 
then
    echo "RELEASING PORTS..."
    pid="$(sudo lsof -i :$pport | awk '{print $2}' | head -2 | tail -1)"
    while [[ $pid -gt 0 ]]
    do
        sudo kill -9 $pid
        pid="$(sudo lsof -i :$pport | awk '{print $2}' | head -2 | tail -1)"
    done

    pid="$(sudo lsof -i :$nport | awk '{print $2}' | head -2 | tail -1)"
    while [[ $pid -gt 0 ]]
    do
        echo "$(expr length "$pid")"
        sudo kill -9 $pid
        pid="$(sudo lsof -i :$nport | awk '{print $2}' | head -2 | tail -1)"
    done

    pid="$(sudo lsof -i :$sport | awk '{print $2}' | head -2 | tail -1)"
    while [[ $pid -gt 0 ]]
    do
        echo "$(expr length "$pid")"
        sudo kill -9 $pid
        pid="$(sudo lsof -i :$sport | awk '{print $2}' | head -2 | tail -1)"
    done
    echo "DONE"

else
    echo "Releasing cancelled"
fi