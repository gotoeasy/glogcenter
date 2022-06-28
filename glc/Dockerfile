# ----------------------------
# 编译
# ----------------------------
FROM golang:alpine AS builder

WORKDIR /goproject
COPY . /goproject

ENV GOPROXY https://goproxy.cn,direct
# 编译，关闭CGO(重要)
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build


# ----------------------------
# 打包镜像
# ----------------------------
FROM alpine AS runner
MAINTAINER gotoeasy <gotoeasy@163.com>

# 复制编译好的glc可执行文件
COPY --from=builder /goproject/glc /usr/local/bin/

# 安装bash以便能进入容器
RUN alpine_version=`cat /etc/issue | head -1 | awk '{print $5}'` \
    && echo "https://mirrors.aliyun.com/alpine/v${alpine_version}/main/" > /etc/apk/repositories \
    && apk update && apk upgrade && apk add --no-cache bash bash-doc bash-completion \
    && rm -rf /var/cache/apk/*

EXPOSE 8080
ENTRYPOINT ["glc"]
