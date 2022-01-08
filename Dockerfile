FROM golang:1.17.5-alpine  as build

LABEL maintainer="longfei6671@163.com"

RUN apk add  --update-cache  libc-dev git gcc

WORKDIR /go/src/app

ADD . /go/src/app/fink-download/

RUN cd fink-download && export GOPROXY="https://goproxy.cn,direct" && go mod download && go build -o finker main.go

FROM alpine:latest

LABEL maintainer="longfei6671@163.com"

COPY --from=build /go/src/app/fink-download/finker /var/www/finker/

RUN  mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN chmod +x /var/www/finker/finker

WORKDIR /var/www/finker/

EXPOSE 9081

CMD ["/var/www/finker/finker"]

