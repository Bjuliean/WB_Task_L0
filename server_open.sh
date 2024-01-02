#!/bin/bash

source .env

host=$SERVER_HOST
port=$SERVER_PORTS

xdg-open http://$host:$port