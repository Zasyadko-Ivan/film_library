#!/bin/sh

echo 'name_image="film_library"
name_docker="film"
workdir="/film_library"
pwd='"$(pwd)" > .env

docker rm $(docker ps -a -q)
docker rmi $(docker images -a -q) 

docker-compose up --build
