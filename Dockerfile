FROM golang:1.18

RUN sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list

RUN apt-get -q update && apt-get -qy install netcat

RUN go env -w GO111MODULE=on

RUN go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER "songcser"

WORKDIR /data

ADD . /data

CMD go mod tidy

CMD go get -u github.com/gin-gonic/gin

RUN go build cmd/main.go

EXPOSE 8080

# ENTRYPOINT ["./main"]
