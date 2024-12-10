package main

import (
	"fmt"
	"runtime"
)

func main() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Allocated Memory: %.2f MB\n", float64(memStats.Alloc)/1024/1024)
}
