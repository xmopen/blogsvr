# TODO: 先从docker容器内拷贝出bin
# 停止容器
docker stop blogsvr
# 删除容器,按照name删除
docker rm blogsvr
# 删除镜像
docker rmi blogsvr:latest

docker build -t blogsvr .

go mod tidy

go build -o ./svrmain ./*.go
# docker run -e TZ=Asia/Shanghai
docker run -d -p 8848:8848 -e TZ=Asia/Shanghai --memory=50m --cpus=0.3 --oom-kill-disable=true  --name blogsvr -v /data/config:/data/config blogsvr:latest