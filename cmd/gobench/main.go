package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/zaluty/gobench/internal/runner"
	"github.com/zaluty/gobench/internal/ui"
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

	fmt.Println(" 🚀 GoBench - Modern Go Benchmarking")
	fmt.Printf("Searching for benchmarks in: %v\n", config.Dirs)
	if config.Dirs == nil {
		fmt.Println("No directories specified. Make sure your benchmarks are in the `benchmarks` directory.")
	}

	// Start in watch mode if specified
	if config.Watch {
		fmt.Println("  Watching for changes...")

		// Initialize terminal UI
		termUI, err := ui.NewTerminalUI()
		if err != nil {
			log.Fatal(err)
		}
		termUI.Start()
		defer termUI.Stop()

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
			fileName := strings.Split(event.Path, "/")[len(strings.Split(event.Path, "/"))-1]
			termUI.AddBenchmark(fileName)
			termUI.SetStatus(fileName, ui.StatusRunning)

			startTime := time.Now()
			cmd := exec.Command("go", "run", event.Path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				termUI.Failed(fileName, err)
				continue
			}

			duration := time.Since(startTime).Seconds()
			termUI.UpdateResult(fileName, duration)
			termUI.Complete(fileName)
		}
	} else {
		// Run benchmarks without watch mode
		runner, err := runner.NewRunner(config.Dirs, config.Filter, config.Parallel)
		if err != nil {
			log.Fatal(err)
		}

		if err := runner.Run(); err != nil {
			log.Fatal(err)
		}
	}
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
