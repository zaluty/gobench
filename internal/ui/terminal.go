package ui

import (
	"fmt"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type BenchmarkStatus string

const (
	StatusPending  BenchmarkStatus = "PENDING"
	StatusRunning  BenchmarkStatus = "RUNNING"
	StatusPassed   BenchmarkStatus = "PASSED"
	StatusFailed   BenchmarkStatus = "FAILED"
	StatusComplete BenchmarkStatus = "COMPLETE"
)

type BenchmarkResult struct {
	Name       string
	Status     BenchmarkStatus
	Duration   float64
	Operations int64
	Memory     float64
	Results    []float64
}

type TerminalUI struct {
	header    *widgets.Paragraph
	table     *widgets.Table
	progress  *widgets.Gauge
	stats     *widgets.Paragraph
	grid      *termui.Grid
	results   map[string]*BenchmarkResult
	startTime time.Time
	running   bool
}

func NewTerminalUI() (*TerminalUI, error) {
	if err := termui.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize termui: %w", err)
	}

	header := widgets.NewParagraph()
	header.Title = "GoBench"
	header.Text = "Modern Benchmarking Tool for Go"
	header.BorderStyle.Fg = termui.ColorCyan

	table := widgets.NewTable()
	table.Title = "Benchmarks"
	table.Rows = [][]string{
		{"Name", "Status", "Duration (s)", "Ops/sec", "Memory (MB)"},
	}
	table.TextStyle = termui.NewStyle(termui.ColorWhite)
	table.TextAlignment = termui.AlignCenter
	table.BorderStyle.Fg = termui.ColorWhite
	table.RowSeparator = true
	table.FillRow = true
	table.RowStyles[0] = termui.NewStyle(termui.ColorYellow, termui.ColorClear, termui.ModifierBold)

	progress := widgets.NewGauge()
	progress.Title = "Progress"
	progress.Percent = 0
	progress.BarColor = termui.ColorGreen
	progress.BorderStyle.Fg = termui.ColorWhite

	stats := widgets.NewParagraph()
	stats.Title = "Summary"
	stats.BorderStyle.Fg = termui.ColorWhite

	ui := &TerminalUI{
		header:   header,
		table:    table,
		progress: progress,
		stats:    stats,
		results:  make(map[string]*BenchmarkResult),
	}

	ui.grid = termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	ui.grid.SetRect(0, 0, termWidth, termHeight)

	ui.grid.Set(
		termui.NewRow(0.1, termui.NewCol(1.0, header)),
		termui.NewRow(0.6, termui.NewCol(1.0, table)),
		termui.NewRow(0.15, termui.NewCol(1.0, progress)),
		termui.NewRow(0.15, termui.NewCol(1.0, stats)),
	)

	return ui, nil
}

func (ui *TerminalUI) AddBenchmark(name string) {
	ui.results[name] = &BenchmarkResult{
		Name:    name,
		Status:  StatusPending,
		Results: make([]float64, 0),
	}
	ui.updateTable()
}

func (ui *TerminalUI) UpdateResult(name string, duration float64) {
	result, exists := ui.results[name]
	if !exists {
		ui.AddBenchmark(name)
		result = ui.results[name]
	}

	result.Status = StatusRunning
	result.Duration = duration
	result.Results = append(result.Results, duration)

	// Calculate operations per second (example metric)
	if duration > 0 {
		result.Operations = int64(1000000000 / duration) // nanoseconds to ops/sec
	}

	ui.updateProgress()
	ui.updateTable()
	ui.updateStats()
	ui.render()
}

func (ui *TerminalUI) updateTable() {
	headers := []string{"Name", "Status", "Duration (s)", "Ops/sec", "Memory (MB)"}
	rows := [][]string{headers}
	for _, result := range ui.results {
		rows = append(rows, []string{
			result.Name,
			string(result.Status),
			fmt.Sprintf("%.2f", result.Duration),
			fmt.Sprintf("%d", result.Operations),
			fmt.Sprintf("%.2f", result.Memory),
		})
	}
	ui.table.Rows = rows
	ui.table.RowStyles[0] = termui.NewStyle(termui.ColorYellow, termui.ColorClear, termui.ModifierBold)
}

func (ui *TerminalUI) updateProgress() {
	if len(ui.results) == 0 {
		ui.progress.Percent = 0
		return
	}

	total := len(ui.results)
	completed := 0
	for _, result := range ui.results {
		if result.Status == StatusComplete || result.Status == StatusFailed {
			completed++
		}
	}
	ui.progress.Percent = int((float64(completed) / float64(total)) * 100)
}

func (ui *TerminalUI) updateStats() {
	// Example stats calculation
	totalDuration := 0.0
	totalOperations := int64(0)
	for _, result := range ui.results {
		totalDuration += result.Duration
		totalOperations += result.Operations
	}
	ui.stats.Text = fmt.Sprintf("Total Duration: %.2f, Total Operations: %d", totalDuration, totalOperations)
}

func (ui *TerminalUI) render() {
	termui.Render(ui.grid)
}

func (ui *TerminalUI) Start() {
	ui.running = true
	ui.startTime = time.Now()
	ui.render()

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		uiEvents := termui.PollEvents()
		for ui.running {
			select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>":
					ui.Stop()
					return
				case "<Resize>":
					payload := e.Payload.(termui.Resize)
					ui.grid.SetRect(0, 0, payload.Width, payload.Height)
					ui.render()
				}
			case <-ticker.C:
				ui.render()
			}
		}
	}()
}

func (ui *TerminalUI) Stop() {
	ui.running = false
	termui.Close()
}

func (ui *TerminalUI) SetStatus(name string, status BenchmarkStatus) {
	if result, exists := ui.results[name]; exists {
		result.Status = status
		ui.updateTable()
		ui.render()
	}
}

func (ui *TerminalUI) Complete(name string) {
	if result, exists := ui.results[name]; exists {
		result.Status = StatusComplete
		ui.updateProgress()
		ui.updateTable()
		ui.updateStats()
		ui.render()
	}
}

func (ui *TerminalUI) Failed(name string, err error) {
	if result, exists := ui.results[name]; exists {
		result.Status = StatusFailed
		ui.updateProgress()
		ui.updateTable()
		ui.updateStats()
		ui.render()
	}
}
