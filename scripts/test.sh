#!/bin/bash

docker-compose up --build

while true; do
    for i in {1..100}; do 
        curl "http://localhost:8080?client_id=$i"
        sleep 0.1
    done
done
