#!/bin/bash
set -e

git pull 

go vet ./...

GOOS=linux go build -o tv server.go

CONTAINER="theoktv"
docker stop $CONTAINER || true
docker rm $CONTAINER ||true

docker run -d --name $CONTAINER --restart always \
	--memory 100m \
	-p "127.0.0.1:9401:8080" \
	-v /etc/ssl:/etc/ssl:ro \
	-v /etc/localtime:/etc/locatime:ro \
	-v $(pwd):/tv \
	debian:jessie /tv/tv -dir /tv

docker logs -f $CONTAINER

