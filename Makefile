.PHONY: help setup test coverage bench lint clean docker-up docker-down integration-test check

help:
	@echo "Go-SQLBuilder Development Commands"
	@echo "==================================="
	@echo "make setup              - 初始化开发环境"
	@echo "make test               - 运行单元测试 (-race)"
	@echo "make coverage           - 生成覆盖率报告"
	@echo "make bench              - 运行基准测试"
	@echo "make lint               - 代码质量检查"
	@echo "make docker-up          - 启动测试数据库"
	@echo "make docker-down        - 停止测试数据库"
	@echo "make integration-test   - 运行集成测试"
	@echo "make check              - 完整检查 (lint + test + coverage)"
	@echo "make clean              - 清理临时文件"

# 初始化开发环境
setup:
	@echo "Setting up development environment..."
	go mod download
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install gotest.tools/gotestsum@latest
	@echo "Setup complete!"

# 单元测试 (带竞态检测)
test:
	@echo "Running unit tests with race detection..."
	go test ./... -v -race -count=1 -timeout=5m

# 生成覆盖率报告
coverage:
	@echo "Generating coverage report..."
	go test ./... -coverprofile=coverage.out -covermode=atomic -timeout=5m
	go tool cover -func=coverage.out | tail -1
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage HTML: coverage.html"

# 基准测试
bench:
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem -benchtime=5s -run=^$$ ./...

# 代码质量检查
lint:
	@echo "Running lint checks..."
	golangci-lint run ./... --timeout=5m

# 启动Docker容器
docker-up:
	@echo "Starting test databases..."
	docker-compose up -d
	@echo "Waiting for databases to be ready..."
	sleep 5
	docker-compose ps

# 停止Docker容器
docker-down:
	@echo "Stopping test databases..."
	docker-compose down

# 集成测试 (需要Docker)
integration-test: docker-up
	@echo "Running integration tests..."
	@sleep 3
	go test -v -tags=integration -timeout=10m ./...
	@echo "Integration tests complete"
	$(MAKE) docker-down

# 性能分析 (CPU profiling)
profile-cpu:
	@echo "Running CPU profile..."
	go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. -benchtime=10s ./...
	go tool pprof -text cpu.prof | head -20
	@echo "Profiles saved: cpu.prof, mem.prof"

# 内存泄漏检测
profile-mem:
	@echo "Analyzing memory allocations..."
	go tool pprof -alloc_space -text mem.prof | head -20

# 完整检查
check: lint test coverage
	@echo "✅ All checks passed!"

# 清理临时文件
clean:
	@echo "Cleaning up..."
	rm -f coverage.out coverage.html
	rm -f cpu.prof mem.prof
	rm -f *.test
	go clean ./...
	@echo "Cleanup complete"

# 安装依赖
install-deps:
	@echo "Installing dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies installed"

# 更新依赖
update-deps:
	@echo "Updating dependencies..."
	go get -u -t ./...
	go mod tidy
	@echo "Dependencies updated"

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# 代码标准化
vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "Go vet passed"

# 快速测试 (仅单元测试)
quick-test:
	@echo "Running quick tests..."
	go test ./... -short -timeout=2m

# 详细的性能基准
detailed-bench:
	@echo "Running detailed benchmarks..."
	go test -bench=. -benchmem -benchtime=10s -run=^$$ -v ./... > bench_results.txt
	@echo "Results saved to: bench_results.txt"
	@cat bench_results.txt | tail -20

# 检查依赖漏洞
check-security:
	@echo "Checking for security vulnerabilities..."
	go list -json -m all | nancy sleuth

# 文档生成
docs:
	@echo "Generating documentation..."
	@echo "API documentation is embedded in code comments"
	@echo "See: README.md, ARCHITECTURE.md, MODERNIZATION_PLAN.md"

# 版本信息
version:
	@echo "Go-SQLBuilder Version Info"
	@echo "=========================="
	@go version
	@go list -m github.com/kamalyes/go-sqlbuilder 2>/dev/null || echo "Current module: github.com/kamalyes/go-sqlbuilder"

# 持续集成 (CI环境)
ci: lint test coverage
	@echo "✅ CI checks passed!"

# 全量构建和测试 (完整流程)
all: clean install-deps fmt vet check docs
	@echo "✅ Full build and test complete!"
