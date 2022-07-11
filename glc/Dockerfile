# ----------------------------
# 编译Web
# ----------------------------
FROM node:14-alpine3.15 AS webbuilder

WORKDIR /webproject
COPY . /webproject

RUN cd www/web && \
    yarn install --registry https://registry.npm.taobao.org && \
    yarn run build


# ----------------------------
# 编译Go
# ----------------------------
FROM golang:alpine AS gobuilder

WORKDIR /goproject
COPY . /goproject

COPY --from=webbuilder /webproject/www/web/dist www/web/

ENV GOPROXY https://goproxy.cn,direct
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build


# ----------------------------
# 打包镜像
# ----------------------------
FROM alpine AS runner
MAINTAINER gotoeasy <gotoeasy@163.com>

COPY --from=gobuilder /goproject/glc /usr/local/bin/

# 安装bash、设定默认时区
RUN alpine_version=`cat /etc/issue | head -1 | awk '{print $5}'` \
    && echo "https://mirrors.aliyun.com/alpine/v${alpine_version}/main/" > /etc/apk/repositories \
    && apk update && apk upgrade && apk add --no-cache bash bash-doc bash-completion tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone \
    && rm -rf /var/cache/apk/*

EXPOSE 8080
ENTRYPOINT ["glc"]
