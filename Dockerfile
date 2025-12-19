# 编译源码
FROM golang:1.25-alpine AS build
COPY . /go/src/example
WORKDIR /go/src/example
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /example ./cmd/

# 构建最终镜像
FROM alpine:latest
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

COPY --from=build /example /app/example
COPY --from=build /go/src/example/config/ /app/config/

EXPOSE 8089
WORKDIR /app
ENTRYPOINT ["./example","-config","./config"]