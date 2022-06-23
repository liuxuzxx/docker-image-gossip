#!/bin/sh
###
 # @Description:  
### 
###
 # @Description:  增加自动构建镜像和推送镜像的shell脚本
### 
cd cmd/gossiper/
go build
cd ../../
docker build -f ./Dockerfile -t liuxuzxx/docker-image-gossip:v1.0.0 .