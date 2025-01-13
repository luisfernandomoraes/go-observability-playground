# Golang Observability Playground

This project is a playground to explore observability practices in Go applications. It includes examples of tracing, metrics, and logging using OpenTelemetry, as well as profiling and load testing.

## Features

### Data Transformations

- Filter users by age (above 30)
- Sort users by name
- Group users by age range
- Batch update user ages
- Count users above a certain age

### Validations

- Input data validation
- Mandatory field checks
- Data format validation

### Observability

#### Tracing (OpenTelemetry)

- Full HTTP request tracing
- Custom spans for each operation
- Detailed attributes in spans
- Context propagation
- Integration with OTLP exporter

#### Metrics (OpenTelemetry)

- Total HTTP requests
- Request duration
- Custom counters per operation
- Latency histograms
- Labels for metric segmentation

#### Logging (slog)

- Structured logs
- Log levels (Info, Warn, Error)
- Context-enriched logs
- Tracing integration

### REST API

- POST endpoint `/users/transform`
- JSON responses
- HTTP error handling
- Middleware for metrics
- Graceful shutdown

### Infrastructure

- Docker containerization
- OpenTelemetry Collector configuration
- Makefile for task automation
- Automated tests

### Profiling

The project includes full support for `pprof` for performance analysis. Available commands include:

#### Profile Collection

- `make pprof-cpu` - Collects a 30-second CPU profile
- `make pprof-mem` - Collects a memory profile
- `make pprof-goroutine` - Analyzes active goroutines
- `make pprof-thread` - Analyzes thread creation
- `make pprof-block` - Analyzes goroutine blocking
- `make trace` - Generates a 30-second execution trace

#### Web Visualization

- `make pprof-web-cpu` - View CPU profile in the browser
- `make pprof-web-mem` - View memory profile in the browser

#### Benchmarking with Profiling

- `make bench-profile` - Runs benchmarks with profiling
- `make bench-cpu` - Views CPU profile from benchmarks
- `make bench-mem` - Views memory profile from benchmarks

#### Cleanup

- `make clean-profiles` - Removes generated profile files

Access profiles via browser at:
- http://localhost:6060/debug/pprof/

### Load Testing with Profiling

The project supports integrated profiling with load testing using `k6` and `pprof`.

#### Commands

- `make profile-load-test` - Runs a load test with full profiling
- `make analyze-load-test` - Analyzes load test results via CLI
- `make view-load-test` - Views load test results in the browser
- `make clean-load-test` - Removes generated profile files

#### Collected Profiles

- CPU Profile: CPU usage analysis during the test
- Heap Profile: Memory usage comparison before and after the test
- Goroutine Profile: Analysis of active goroutines
- Execution Trace: Detailed execution trace

#### Visualization Ports

- CPU Profile: http://localhost:8081
- Heap Profile: http://localhost:8082
- Execution Trace: Default browser

## Requirements

- Go 1.21+
- Docker and Docker Compose
- OpenTelemetry Collector

## How to Run

1. Install dependencies:

```bash
   make install-tools
```

### Build, Test, Check Vulnerabilities, Lint and Generate Code Coverage

```bash
  make
```

### Run the OpenTelemetry Stack

```bash
  make run-otel-stack
```

### Run the Application

```bash
  make run
```

### Run Load Test

```bash
  make load-test
```

### Jaeger UI

[Jaeger](http://localhost:16686/>)

Thanks for your time! 😊

#### TODO

- [ ] ELK Stack
- [-] OpenTelemetry
- [ ] Datadog
- [ ] Prometheus & Grafana
- [ ] Dockerize
- [ ] Kubernetes
- [ ] GitHub Actions
- [ ] Swagger
- [ ] API Gateway
