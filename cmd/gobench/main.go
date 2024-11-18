package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zaluty/gobench/internal/runner"
)

// CLI flags
type Config struct {
	Watch    bool     // Watch mode
	Filter   string   // Benchmark filter
	Dirs     []string // Directories to scan
	Reporter string   // Reporter type (terminal, json, csv)
	Parallel int      // Number of parallel runners
}

func main() {
	config := parseFlags()

	// TODO: Initialize components
	// - File watcher
	// - Benchmark runner
	// - Reporter
	// - Analyzer

	fmt.Println(" 🚀 GoBench - Modern Go Benchmarking")
	fmt.Printf("Searching for benchmarks in: %v\n", config.Dirs)
	if config.Dirs == nil {
		fmt.Println("No directories specified. Make sure your benchmarks ar in the `benchmarks` directory.")
	}
	// Start in watch mode if specified
	if config.Watch {
		fmt.Println("  Watching for changes...")

		watcher, err := runner.NewWatcher(config.Dirs)
		if err != nil {
			log.Fatal(err)
		}

		err = watcher.Start()
		if err != nil {
			log.Fatal(err)
		}

		// Watch for events and execute files
		for event := range watcher.Events() {
			fmt.Printf("\nExecuting file: %s\n", event.Path)
			cmd := exec.Command("go", "run", event.Path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error executing %s: %v\n", event.Path, err)
			}
		}
	}

	// TODO: Discover and run benchmarks
}

func parseFlags() *Config {
	config := &Config{}

	flag.BoolVar(&config.Watch, "watch", false, "Watch mode")
	flag.StringVar(&config.Filter, "filter", "", "Benchmark filter")
	flag.StringVar(&config.Reporter, "reporter", "terminal", "Reporter type (terminal, json, csv)")
	flag.IntVar(&config.Parallel, "parallel", 1, "Number of parallel runners")

	flag.Parse()

	// Use remaining args as directories to scan, default to current directory
	args := flag.Args()
	if len(args) == 0 {
		defaultDir := "./benchmarks"
		absPath, err := filepath.Abs(defaultDir)
		if err != nil {
			log.Fatal(err)
		}
		config.Dirs = []string{absPath}
	} else {
		config.Dirs = make([]string, len(args))
		for i, arg := range args {
			absPath, err := filepath.Abs(arg)
			if err != nil {
				log.Fatal(err)
			}
			config.Dirs[i] = absPath
		}
	}

	return config
}
