# 1、删除老的Deployment 对象和 minikube 缓存镜像
kubectl delete deployment blogsvr
minikube cache delete openxm/blogsvr:latest

# 2、 删除 minikube 内老的原始镜像
eval $(minikube docker-env)
# 强制删除 minikube 内对应的镜像，上层已经将缓存内的镜像也已删除掉
docker rmi -f openxm/blogsvr:latest
eval $(minikube docker-env -u)

# 3、重新打包二进制和镜像
go mod tidy
go build -o ./svrmain ./*.go

# 4、删除宿主机镜像并重新构建
docker rmi openxm/blogsvr:latest
docker build -t openxm/blogsvr:latest .

# 5、添加到 minikube 缓存内
minikube cache add openxm/blogsvr:latest
# 6、创建 Deployment 对象
kubectl apply -f ./blogsvr_pod.yaml

