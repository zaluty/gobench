package main

import (
	"fmt"
	"runtime"
	"time"
)

const (
	// Allocate 1GB of memory in chunks
	totalMemory = 1024 * 1024 * 1024 // 1 GB
	chunkSize   = 1024 * 1024        // 1 MB chunks
)

func allocateMemory() [][]byte {
	chunks := make([][]byte, 0, totalMemory/chunkSize)

	fmt.Println("Starting memory allocation...")
	for i := 0; i < totalMemory/chunkSize; i++ {
		// Allocate 1MB chunks
		chunk := make([]byte, chunkSize)
		// Write some data to ensure it's actually allocated
		for j := 0; j < len(chunk); j += 1024 {
			chunk[j] = byte(i % 256)
		}
		chunks = append(chunks, chunk)

		if i%100 == 0 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Allocated: %d MB, System Memory: %d MB\n",
				i+1, m.Sys/1024/1024)
		}
	}
	return chunks
}

func main() {
	fmt.Println("Starting memory benchmark...")

	startTime := time.Now()
	chunks := allocateMemory()

	// Force GC to get accurate memory stats
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("\nBenchmark Results:\n")
	fmt.Printf("Time taken: %v\n", time.Since(startTime))
	fmt.Printf("Total memory allocated: %d MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("System memory in use: %d MB\n", m.Sys/1024/1024)
	fmt.Printf("Number of chunks: %d\n", len(chunks))

	// Keep the memory allocated for a moment to see it in the UI
	time.Sleep(2 * time.Second)
}
