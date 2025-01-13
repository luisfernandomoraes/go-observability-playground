#!/bin/bash

# Get the project's root directory (one level above the scripts directory)
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# Create a directory for the profiles if it doesn't exist
mkdir -p "$PROJECT_ROOT/profiles"

# Start collecting the CPU profile
echo "Starting CPU profile..."
curl -o "$PROJECT_ROOT/profiles/cpu.prof" http://localhost:6060/debug/pprof/profile?seconds=35 &
CPU_PID=$!

# Collect the heap profile (before the test)
echo "Collecting initial heap profile..."
curl -o "$PROJECT_ROOT/profiles/heap-before.prof" http://localhost:6060/debug/pprof/heap

# Start collecting the trace
echo "Starting trace..."
curl -o "$PROJECT_ROOT/profiles/trace.out" http://localhost:6060/debug/pprof/trace?seconds=35 &
TRACE_PID=$!

# Wait 2 seconds to ensure profiling has started
sleep 2

# Run the load test
echo "Running load test..."
docker run --rm -i \
    --add-host=host.docker.internal:host-gateway \
    -v "$PROJECT_ROOT/load-tests:/scripts" \
    -w /scripts \
    grafana/k6 run k6.js

# Wait for the CPU profile and trace to finish
wait $CPU_PID
wait $TRACE_PID

# Collect the heap profile after the test
echo "Collecting final heap profile..."
curl -o "$PROJECT_ROOT/profiles/heap-after.prof" http://localhost:6060/debug/pprof/heap

# Collect the goroutine profile
echo "Collecting goroutine profile..."
curl -o "$PROJECT_ROOT/profiles/goroutine.prof" http://localhost:6060/debug/pprof/goroutine

echo "Profiles collected in $PROJECT_ROOT/profiles/"