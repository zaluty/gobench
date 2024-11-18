# GoBench - Modern Benchmarking Tool for Go

A fast, modern benchmarking tool for Go developers inspired by Vitest.

## Project Structure

```
.
├── cmd/
│   └── gobench/           # Main entry point
│       └── main.go
├── internal/
│   ├── runner/           # Benchmark runner implementation
│   │   ├── runner.go     # Core benchmark runner
│   │   ├── watcher.go    # File watcher
│   │   └── parser.go     # Benchmark file parser
│   ├── reporter/         # Results reporting
│   │   ├── terminal.go   # Terminal UI
│   │   ├── formatter.go  # Output formatting
│   │   └── export.go     # Export results
│   ├── analyzer/         # Benchmark analysis
│   │   ├── compare.go    # Compare results
│   │   └── stats.go      # Statistical analysis
│   └── config/           # Configuration
│       └── config.go     # Config structs and loading
├── pkg/
│   └── utils/            # Shared utilities
│       ├── colors.go     # Terminal colors
│       └── format.go     # String formatting
├── go.mod
└── README.md
```

## Core Components

1. **Runner**: 
   - Discovers and executes benchmarks
   - Watches for file changes
   - Manages parallel execution

2. **Reporter**:
   - Beautiful terminal UI
   - Real-time results display
   - Export functionality

3. **Analyzer**:
   - Statistical analysis
   - Performance comparisons
   - Memory profiling

4. **Configuration**:
   - CLI flags
   - Config file support
   - Benchmark filters

## Key Features

- 🔄 Hot reloading
- 📊 Real-time performance graphs
- 🚀 Parallel execution
- 📈 Historical comparisons
- 🎨 Beautiful terminal UI
- 📤 Export results to JSON/CSV
