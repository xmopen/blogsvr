# scratch
FROM centos
LABEL authors="openxm"

ENV TZ=Asia/Shanghai

# 临时方案：
# 1. 拷贝data目录
#COPY /data/ /data/
COPY ./svrmain /data/code/go/blogsvr/
COPY ./config/ /data/config/

EXPOSE 8848
ENTRYPOINT ["/data/code/go/blogsvr/svrmain"]