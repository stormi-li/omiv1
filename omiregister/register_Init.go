package register

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/process"
)

func Init(r *Register) {
	r.AddRegisterHandleFunc("Weight", func() string {
		return strconv.Itoa(r.Weight)
	})
	pid := os.Getpid()
	r.AddRegisterHandleFunc("Pid", func() string {
		return strconv.Itoa(pid)
	})
	host, _ := os.Hostname()
	r.AddRegisterHandleFunc("Host", func() string {
		return host
	})
	startTime := r.StartTime.Format("2006-01-02 15:04:05")
	r.AddRegisterHandleFunc("StartTime", func() string {
		return startTime
	})
	r.AddRegisterHandleFunc("RunTime", func() string {
		elapsed := time.Since(r.StartTime)
		seconds := int64(elapsed.Seconds())
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		secs := seconds % 60
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	})

	virtualMem, _ := mem.VirtualMemory()
	memTotal := float64(virtualMem.Total)

	p, _ := process.NewProcess(int32(pid))
	r.AddRegisterHandleFunc("Memory", func() string {
		memInfo, _ := p.MemoryInfo()
		memPercent := (float64(memInfo.RSS) / (memTotal + 1)) * 100
		return fmt.Sprintf("%dMB %.2f%%", memInfo.RSS/1024/1024, memPercent)
	})
	r.AddRegisterHandleFunc("CPU", func() string {
		cpuPercent, _ := p.CPUPercent()
		return fmt.Sprintf("%.2f%%", cpuPercent)
	})

	r.AddRegisterHandleFunc("MessageHandlers", func() string {
		handlerNames := []string{}
		for name := range r.MessageHandler.handleFuncs {
			handlerNames = append(handlerNames, name)
		}
		sort.Slice(handlerNames, func(i, j int) bool {
			return handlerNames[i] < handlerNames[j]
		})
		return strings.Join(handlerNames, ", ")
	})

	r.AddMessageHandleFunc(Command_UpdateWeight, func(message string) {
		weight, err := strconv.Atoi(message)
		if err == nil {
			r.Weight = weight
		}
	})
}
