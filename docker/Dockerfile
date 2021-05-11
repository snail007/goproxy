# alpine, busybox, scratch
FROM scratch
COPY proxy /
COPY Shanghai /etc/localtime/
COPY ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/proxy"]
CMD ["--version"]