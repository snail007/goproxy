FROM golang:1.10.3-alpine
ARG GOPROXY_VERSION=master
RUN apk update; apk upgrade; \
    apk add --no-cache git; \
    cd /go/src/; \
    mkdir github.com; \
    mkdir github.com/snail007; \
    cd github.com/snail007; \
    git clone https://github.com/snail007/goproxy.git; \
	cd goproxy; \
    git checkout ${GOPROXY_VERSION}; \
    go build -ldflags "-s -w" -o proxy; \
    cp proxy /proxy
FROM alpine:3.7
CMD /proxy ${OPTS}