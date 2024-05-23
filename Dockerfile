FROM golang:1.22-alpine AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache build-base
COPY . /app
RUN ls /app
WORKDIR /app

ENV GOPROXY https://goproxy.cn,direct
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -o /go/bin/app_release


FROM alpine

RUN mkdir -p /opt/conf
COPY --from=builder /go/bin/app_release /opt/app

COPY ./conf/config.yaml /opt/conf/config.yaml
WORKDIR /opt

CMD ["/opt/app"]
