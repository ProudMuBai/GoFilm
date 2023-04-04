FROM golang:1.20

ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
        GOPROXY=https://goproxy.cn,direct

MAINTAINER "bai"
 
WORKDIR /opt/server
 
ADD ./server /opt/server
 
RUN go build main.go
 
EXPOSE 3061