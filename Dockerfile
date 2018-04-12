FROM golang:1.10.1 as builder
RUN mkdir -p $GOPATH/src/github.com/solarwinds/gitlic-check
WORKDIR $GOPATH/src/github.com/solarwinds/gitlic-check
ADD . .
RUN GOOS=linux GARCH=amd64 go build -a -o /bin/gitlic-check

FROM alpine:3.7
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /bin/gitlic-check /bin/gitlic-check
CMD gitlic-check