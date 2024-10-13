.PHONY: all build run gotool clean help
BINARY="statefulset"
all: gotool build
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}
build:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BINARY}.exe
run:
	@go run ./
gotool:
	go fmt ./
	go vet ./
test:
	go test -coverprofile=size_coverage.out ./base/...
clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
help:
	@echo "make build0linux- 格式化 Go 代码, 并编译生成linux二进制文件"
	@echo "make build - 编译 Go 代码, 生成windows二进制文件"
	@echo "make test - 直接运行 Go test覆盖测试内容"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"