BINARY_NAME=example
.PHONY: all dep build-win build-linux
# 默认目标
all: build
# 更新依赖
dep:
	go mod tidy
# 编译
build-win:
	go mod tidy
	go env -w CGO_ENABLED=0
	go env -w GOARCH=amd64
	go env -w GOOS=windows
	go build -o bin/$(BINARY_NAME).exe ./...

build-linux:
	go mod tidy
	go env -w CGO_ENABLED=0
	go env -w GOARCH=amd64
	go env -w GOOS=linux
	go build -o bin/$(BINARY_NAME) ./...	
