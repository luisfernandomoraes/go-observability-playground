.DEFAULT_GOAL := all

.PHONY: all
all: build test vulncheck lint

.PHONY: build
build:
	go build -o bin/app cmd/main.go
	@echo "📦 Build Done"

.PHONY: test
test:
	go test -v -race -cover ./...
	@echo "🧪 Test Completed"

.PHONY: vulncheck
vulncheck:
	@echo "🔍 Running Vulnerability Check"
	govulncheck ./...
	@echo "✅ Vulnerability Check Completed"

.PHONY: lint
lint:
	@echo "🔎 Running Linter"
	golangci-lint run ./...
	@echo "✅ Linter Check Completed"

.PHONY: run
run:
	@echo "🚀 Building Optimized Binary..."
    # CGO_ENABLED=0: Disables CGO (C bindings), forcing the Go compiler to produce a statically linked binary.
	# -s: Strips the symbol table from the binary.
    # -w: Omits the DWARF debugging information.
	# CGO_ENABLED=0 go build -o bin/app -ldflags="-s -w" cmd/main.go
	go build -o bin/app cmd/main.go
	@echo "🚀 Running Application..."
	./bin/app
	@echo "✅ Application Completed"

.PHONY: up-otel-stack
up-otel-stack:
	@echo "🚀 Starting OpenTelemetry Stack"
	docker compose up -d
	@echo "✅ OpenTelemetry Stack"

.PHONY: down-otel-stack
down-otel-stack:
	@echo "🚀 Stopping OpenTelemetry Stack"
	docker compose down
	@echo "✅ OpenTelemetry Stack Stopped"

.PHONY: load-test
load-test:
	@echo "Running Load Test"
	docker run --rm -i --add-host=host.docker.internal:host-gateway -v $(PWD)/load-tests:/scripts -w /scripts grafana/k6 run k6.js
	@echo "Load Test Completed"

.PHONY: clear
clear:
	@echo "🧹 Cleaning"
	rm -rf bin
	@echo "🧹 Cleaned"

.PHONY: install-tools
install-tools:
	@echo "📦 Installing Required Tools"
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/google/pprof@latest
	go mod tidy
	docker pull grafana/k6
	@echo "✅ Dependencies Installed"

.PHONY: pprof-cpu
pprof-cpu:
	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

.PHONY: pprof-mem
pprof-mem:
	go tool pprof http://localhost:6060/debug/pprof/heap

.PHONY: pprof-goroutine
pprof-goroutine:
	go tool pprof http://localhost:6060/debug/pprof/goroutine

.PHONY: pprof-thread
pprof-thread:
	go tool pprof http://localhost:6060/debug/pprof/threadcreate

.PHONY: pprof-block
pprof-block:
	go tool pprof http://localhost:6060/debug/pprof/block

.PHONY: trace
trace:
	curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=30
	go tool trace trace.out

.PHONY: pprof-web-cpu
pprof-web-cpu:
	go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30

.PHONY: pprof-web-mem
pprof-web-mem:
	go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap

.PHONY: bench-profile
bench-profile:
	go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./...

.PHONY: bench-cpu
bench-cpu:
	go tool pprof -http=:8081 cpu.prof

.PHONY: bench-mem
bench-mem:
	go tool pprof -http=:8081 mem.prof

.PHONY: clean-profiles
clean-profiles:
	rm -f *.prof trace.out

.PHONY: profile-load-test
profile-load-test:
	@echo "🔍 Starting profiling with load testing"
	chmod +x scripts/profile-load-test.sh
	./scripts/profile-load-test.sh
	@echo "✅ Profiling completed"

.PHONY: analyze-load-test
analyze-load-test:
	@echo "📊 Analyzing load test results"
	@echo "\n🔍 CPU Profile:"
	go tool pprof -top profiles/cpu.prof
	@echo "\n🔍 Heap Profile (diff):"
	go tool pprof -top -diff_base=profiles/heap-before.prof profiles/heap-after.prof
	@echo "\n🔍 Goroutine Profile:"
	go tool pprof -top profiles/goroutine.prof
	@echo "\n📊 Analysis completed"

.PHONY: view-load-test
view-load-test:
	@echo "🌐 Opening visualizations in the browser"
	go tool pprof -http=:8081 profiles/cpu.prof &
	go tool pprof -http=:8082 profiles/heap-after.prof &
	go tool trace profiles/trace.out &
	@echo "✅ Visualizations opened"

.PHONY: clean-load-test
clean-load-test:
	rm -rf profiles/
	@echo "🧹 Profiles removed"