FROM golang:1.26-alpine AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache build-base

WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /go/bin/app .


FROM alpine:1.22

RUN mkdir -p /opt/conf
COPY --from=builder /go/bin/app /opt/app

WORKDIR /opt
EXPOSE 8080
CMD ["/opt/app"]
