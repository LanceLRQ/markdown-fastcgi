#!/bin/sh
# docker pull nginx:latest


WORK_DIR=$(pwd)
export GOOS=linux
export GOARCH=amd64
go build -o main main.go

CID=`docker run -d -p 8848:80 -v $WORK_DIR/main:/markdown -v $WORK_DIR/test.md:/usr/share/nginx/html/test.md -v $WORK_DIR/nginx-site.conf:/etc/nginx/conf.d/default.conf nginx:latest`
docker exec -it ${CID} bash
docker kill ${CID} && docker rm ${CID}