#!/bin/bash

cp -r ../karate/karate.jar .
cp -r ../karate/log-config.xml .
docker build -t karate .
rm karate.jar
rm log-config.xml

docker run --network host --env ENVIRONMENT=$1 karate /karate/cases/.