# go-health-checker

# Concurrent HTTP Health Checker CLI in Go

A fast, concurrent command-line utility built in Go to check the availability, HTTP status codes, and response latencies of target URLs using worker pools and context deadlines.

---

## Features

- **Concurrent Execution:** Utilizes Go goroutines and buffered channels for parallel URL health checks.
- **Worker Pool Architecture:** Restricts system resource usage by maintaining a controlled pool of active worker goroutines.
- **Context Timeouts:** Enforces request deadlines via `context.WithTimeout` to prevent hanging HTTP connections.
- **Resource Management:** Ensures HTTP response bodies are properly closed to eliminate network socket leaks.

---

## Architecture

```text
                  ┌─> Worker 1 (http.Get) ─┐
URLs -> Channel ->├─> Worker 2 (http.Get) ─┼─> Results Channel -> Terminal Output
                  └─> Worker 3 (http.Get) ─┘

Getting Started
Prerequisites
1.Go installed on your system.

Running Locally
1.Clone the repository:git clone [https://github.com/YOUR_USERNAME/go-health-checker.git](https://github.com/YOUR_USERNAME/go-health-checker.git)
cd go-health-checker
