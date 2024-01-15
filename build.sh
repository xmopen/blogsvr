#docker stop blogsvr
#docker rm blogsvr
#docker rmi blogsvr
#go mod tidy
#go build -o ./svrmain ./*.go
#docker build -t zhenxinma/blogsvr .

# 将镜像添加到minikube中
#minikube cache add zhenxinma/blogsvr:latest

#kubectl delete deployment blogsvr
#kubectl apply -f ./blogsvr_pod.yaml

#docker run -d -p 8848:8848 -e TZ=Asia/Shanghai --memory=50m --cpus=0.3 --oom-kill-disable=true  --name blogsvr -v /data/code/go/blogsvr:/data/app/go/blogsvr/ -v /data/config:/data/config blogsvr:latest



go mod tidy
go build -o ./svrmain ./*.go

docker rmi openxm/blogsvr:latest

docker build -t openxm/blogsvr:latest .
minikube cache delete openxm/blogsvr:latest
minikube cache add openxm/blogsvr:latest

kubectl delete deployment blogsvr
kubectl apply -f ./blogsvr_pod.yaml