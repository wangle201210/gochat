.PHONY: build run clean install help

# 默认目标
all: build

# 构建应用
build:
	@echo "构建 GoChat..."
	@go build -o gochat cmd/gochat/main.go
	@echo "构建完成: ./gochat"

# 运行应用
run:
	@go run cmd/gochat/main.go

# 清理构建产物
clean:
	@echo "清理构建产物..."
	@rm -f gochat
	@rm -f gochat.exe
	@echo "清理完成"

# 安装依赖
install:
	@echo "安装依赖..."
	@go mod tidy
	@go mod download
	@echo "依赖安装完成"

# 跨平台构建
build-all: build-linux build-windows build-darwin

build-linux:
	@echo "构建 Linux 版本..."
	@GOOS=linux GOARCH=amd64 go build -o gochat-linux-amd64 cmd/gochat/main.go

build-windows:
	@echo "构建 Windows 版本..."
	@GOOS=windows GOARCH=amd64 go build -o gochat-windows-amd64.exe cmd/gochat/main.go

build-darwin:
	@echo "构建 macOS 版本..."
	@GOOS=darwin GOARCH=amd64 go build -o gochat-darwin-amd64 cmd/gochat/main.go

# 帮助信息
help:
	@echo "GoChat - Makefile 使用说明"
	@echo ""
	@echo "可用命令:"
	@echo "  make build        - 构建应用"
	@echo "  make run          - 运行应用"
	@echo "  make clean        - 清理构建产物"
	@echo "  make install      - 安装依赖"
	@echo "  make build-all    - 构建所有平台版本"
	@echo "  make build-linux  - 构建 Linux 版本"
	@echo "  make build-windows- 构建 Windows 版本"
	@echo "  make build-darwin - 构建 macOS 版本"
	@echo "  make help         - 显示此帮助信息"
