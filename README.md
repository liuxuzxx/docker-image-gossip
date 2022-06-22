<!--
 * @Description:  描述项目如何使用的文档
-->
# 概述
docker-image-gossip项目主要是为了处理Kubernetes集群当下载k8s.gcr.io/gcr.io开头的Image的时候，会出现拉取不到的问题，现在能想到的一种方式就是：通过docker load 的方式来处理，但是这种方式有如下的问题:
1. 需要登录到任何一个node机器上进行docker load操作
2. 当集群有新的node加入的时候，需要再次手动导入下(可能会出现遗忘的情况)

# 设计
## 功能设计
### 提供导入接口
提供一个HTTP协议的接口，可以导入docker image是tar文件，然后调用本node的docker命令执行:docker load < xxx.tar

# 构建

## 从源码构建
```bash
bash build-docker.sh
```
本脚本负责编译golang源码，并且打包镜像上传到本机的harbor服务和华为云服务！

# 使用

## 使用curl访问
```bash
curl -X POST http://部署的服务ip:port/v1/docker/gossip-image -F "file=@/home/liuxu/Videos/golang-15.tar" -H "Content-Type: multipart/form-data"
```