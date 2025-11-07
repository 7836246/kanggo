.PHONY: build build-fasthttp test bench bench-fasthttp clean help install-deps run-demo run-demo-fasthttp

# 默认目标
.DEFAULT_GOAL := help

# 变量
BINARY_NAME=kanggo-app
BUILD_DIR=bin
GO=go
GOFLAGS=
TAGS=

## help: 显示帮助信息
help:
	@echo "KangGo 构建命令"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "构建命令："
	@echo "  make build             - 构建（net/http 引擎）"
	@echo "  make build-fasthttp    - 构建（fasthttp 引擎）"
	@echo ""
	@echo "测试命令："
	@echo "  make test              - 运行所有测试"
	@echo "  make bench             - 运行基准测试（net/http）"
	@echo "  make bench-fasthttp    - 运行基准测试（fasthttp）"
	@echo ""
	@echo "运行命令："
	@echo "  make run-demo          - 运行演示（net/http）"
	@echo "  make run-demo-fasthttp - 运行演示（fasthttp）"
	@echo ""
	@echo "其他命令："
	@echo "  make install-deps      - 安装依赖"
	@echo "  make clean             - 清理构建文件"
	@echo ""

## build: 构建应用（net/http 引擎）
build:
	@echo "📦 构建中... (net/http 引擎)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)
	@echo "✓ 构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

## build-fasthttp: 构建应用（fasthttp 引擎）
build-fasthttp:
	@echo "⚡ 构建中... (fasthttp 引擎)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -tags fasthttp -o $(BUILD_DIR)/$(BINARY_NAME)-fast
	@echo "✓ 构建完成: $(BUILD_DIR)/$(BINARY_NAME)-fast"

## test: 运行所有测试
test:
	@echo "🧪 运行测试..."
	$(GO) test -v ./...

## bench: 运行基准测试（net/http）
bench:
	@echo "📊 运行基准测试 (net/http)..."
	$(GO) test -bench=. -benchmem -benchtime=2s

## bench-fasthttp: 运行基准测试（fasthttp）
bench-fasthttp:
	@echo "📊 运行基准测试 (fasthttp)..."
	$(GO) test -tags fasthttp -bench=. -benchmem -benchtime=2s

## bench-compare: 对比两种引擎的性能
bench-compare:
	@echo "📊 性能对比测试..."
	@echo ""
	@echo "━━━ net/http 引擎 ━━━"
	@$(GO) test -bench=BenchmarkStaticRoute -benchmem -benchtime=1s | grep "BenchmarkStaticRoute"
	@echo ""
	@echo "━━━ fasthttp 引擎 ━━━"
	@$(GO) test -tags fasthttp -bench=BenchmarkStaticRoute -benchmem -benchtime=1s | grep "BenchmarkStaticRoute" || echo "需要安装 fasthttp: go get github.com/valyala/fasthttp"

## run-demo: 运行演示程序（net/http）
run-demo:
	@echo "🚀 运行演示... (net/http 引擎)"
	$(GO) run examples/phase2_demo.go

## run-demo-fasthttp: 运行演示程序（fasthttp）
run-demo-fasthttp:
	@echo "⚡ 运行演示... (fasthttp 引擎)"
	$(GO) run -tags fasthttp examples/phase2_demo.go

## install-deps: 安装依赖
install-deps:
	@echo "📦 安装依赖..."
	$(GO) get github.com/valyala/fasthttp
	$(GO) mod tidy
	@echo "✓ 依赖安装完成"

## clean: 清理构建文件
clean:
	@echo "🧹 清理中..."
	@rm -rf $(BUILD_DIR)
	@$(GO) clean
	@echo "✓ 清理完成"

## lint: 运行代码检查
lint:
	@echo "🔍 运行代码检查..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "需要安装 golangci-lint"; exit 1; }
	golangci-lint run

## fmt: 格式化代码
fmt:
	@echo "✨ 格式化代码..."
	$(GO) fmt ./...
	@echo "✓ 格式化完成"

## mod: 更新依赖
mod:
	@echo "📦 更新依赖..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✓ 依赖更新完成"
