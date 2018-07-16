FROM golang:1.10.3-alpine
ARG GOPROXY_VERSION=v5.2
RUN apk update; apk upgrade; \
    apk add --no-cache git; \
    cd /go/src/; \
    mkdir github.com; \
    mkdir github.com/snail007; \
    cd github.com/snail007; \
    git clone https://github.com/snail007/goproxy.git; \
	cd goproxy; \
    git checkout ${GOPROXY_VERSION}; \
    go build -ldflags "-s -w" -o proxy;
FROM alpine:3.7
CMD cd /go/src/github.com/snail007/goproxy/ && ./proxy ${OPTS}