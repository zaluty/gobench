package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zaluty/gobench/internal/ui"
)

// BenchmarkResult represents the result of a single benchmark
type BenchmarkResult struct {
	Name     string
	Duration float64
	Ops      int64
	AllocMB  float64
}

// Runner handles benchmark discovery and execution
type Runner struct {
	dirs     []string
	filter   string
	parallel int
	ui       *ui.TerminalUI
}

// NewRunner creates a new benchmark runner
func NewRunner(dirs []string, filter string, parallel int) (*Runner, error) {
	termUI, err := ui.NewTerminalUI()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize UI: %w", err)
	}

	return &Runner{
		dirs:     dirs,
		filter:   filter,
		parallel: parallel,
		ui:       termUI,
	}, nil
}

// Run executes all discovered benchmarks
func (r *Runner) Run() error {
	r.ui.Start()
	defer r.ui.Stop()

	benchmarks, err := r.discover()
	if err != nil {
		return fmt.Errorf("failed to discover benchmarks: %w", err)
	}

	var wg sync.WaitGroup
	resultChan := make(chan BenchmarkResult, len(benchmarks))
	errorChan := make(chan error, len(benchmarks))

	// Initialize UI with discovered benchmarks
	for bench := range benchmarks {
		r.ui.AddBenchmark(bench)
	}

	// Create worker pool
	for i := 0; i < r.parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for bench := range benchmarks {
				r.ui.SetStatus(bench, ui.StatusRunning)

				result, err := r.runBenchmark(bench)
				if err != nil {
					r.ui.Failed(bench, err)
					errorChan <- fmt.Errorf("benchmark %s failed: %w", bench, err)
					continue
				}

				r.ui.UpdateResult(bench, result.Duration)
				r.ui.Complete(bench)
				resultChan <- result
			}
		}()
	}

	// Wait for completion
	wg.Wait()
	close(resultChan)
	close(errorChan)

	// Check for errors
	if len(errorChan) > 0 {
		return <-errorChan
	}

	return nil
}

// discover finds all benchmark functions in the specified directories
func (r *Runner) discover() (chan string, error) {
	benchmarks := make(chan string, 100)

	go func() {
		defer close(benchmarks)
		for _, dir := range r.dirs {
			// Find all *_test.go files
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(path, "_test.go") {
					// TODO: Parse file to find benchmark functions
					// For now, just run go test -bench
					benchmarks <- path
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Error walking directory %s: %v\n", dir, err)
			}
		}
	}()

	return benchmarks, nil
}

// runBenchmark executes a single benchmark
func (r *Runner) runBenchmark(path string) (BenchmarkResult, error) {
	dir := filepath.Dir(path)

	// Run benchmark using go test
	cmd := exec.Command("go", "test", "-bench", r.filter, "-benchmem")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return BenchmarkResult{}, fmt.Errorf("benchmark failed: %w\n%s", err, output)
	}

	// TODO: Parse benchmark output to extract results
	// For now, return dummy result
	return BenchmarkResult{
		Name:     filepath.Base(path),
		Duration: 1.0,
		Ops:      1000000,
		AllocMB:  1.5,
	}, nil
}
