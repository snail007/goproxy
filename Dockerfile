FROM golang:1.10.3-alpine as builder
ARG GOPROXY_VERSION=master
RUN apk update && apk upgrade && \
    apk add --no-cache git && \
    cd /go/src/ && \
    mkdir snail007 && \
    cd snail007 && \
    git clone https://github.com/snail007/goproxy.git && \
	cd goproxy && \
    git checkout ${GOPROXY_VERSION} && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /proxy
FROM alpine:3.7
CMD /proxy ${OPTS}
