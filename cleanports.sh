#!/bin/bash

echo "The script will clear the ports used by the program by default: 8080, 5432, 4222."
echo "This can be useful if the server is not running in a container."
echo "Allow? (y/n)"

read symbol

if [ $symbol = 'y' ] || [ $symbol = 'Y' ] 
then
    echo "RELEASING PORTS..."
    pid="$(sudo lsof -i :5432 | awk '{print $2}' | head -2 | tail -1)"
    while [[ $pid -gt 0 ]]
    do
        sudo kill -9 $pid
        pid="$(sudo lsof -i :5432 | awk '{print $2}' | head -2 | tail -1)"
    done

    pid="$(sudo lsof -i :8080 | awk '{print $2}' | head -2 | tail -1)"
    while [[ $pid -gt 0 ]]
    do
        echo "$(expr length "$pid")"
        sudo kill -9 $pid
        pid="$(sudo lsof -i :8080 | awk '{print $2}' | head -2 | tail -1)"
    done

    pid="$(sudo lsof -i :4222 | awk '{print $2}' | head -2 | tail -1)"
    while [[ $pid -gt 0 ]]
    do
        echo "$(expr length "$pid")"
        sudo kill -9 $pid
        pid="$(sudo lsof -i :4222 | awk '{print $2}' | head -2 | tail -1)"
    done
    echo "DONE"

else
    echo "Releasing cancelled"
fi