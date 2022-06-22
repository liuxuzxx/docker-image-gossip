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
docker build -f ./Dockerfile -t swr.cn-south-1.myhuaweicloud.com/cpaas/component/docker-image-gossip:v1.0.0 .
docker tag swr.cn-south-1.myhuaweicloud.com/cpaas/component/docker-image-gossip:v1.0.0 172.16.15.121:10000/cpaas/component/docker-image-gossip:v1.0.0
docker push 172.16.15.121:10000/cpaas/component/docker-image-gossip:v1.0.0