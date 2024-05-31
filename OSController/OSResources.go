package OSController

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type OSState struct {
	TotalCPUPercent float64
	PerCPUPercent   []float64
	TotalMemPercent float64
	TotalGPUPercent float64
}

func GetOSStates(interval time.Duration) {
	// 获取 CPU 核心数
	cores, err := cpu.Counts(true)
	if err != nil {
		fmt.Println("Failed to get CPU count:", err)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		state, err := collectOSState(cores, interval)
		if err != nil {
			fmt.Println("Failed to collect OS state:", err)
			continue
		}

		// 输出总利用率
		fmt.Printf("Total CPU usage: %.2f%%\n", state.TotalCPUPercent)
		fmt.Println(state.PerCPUPercent)

		// 输出总内存利用率
		fmt.Printf("Total memory usage: %.2f%%\n", state.TotalMemPercent)

		// 输出 GPU 利用率
		fmt.Println("GPU Utilization:", state.TotalGPUPercent, "%")
	}
}

func collectOSState(cores int, interval time.Duration) (*OSState, error) {
	state := &OSState{}

	// 获取每个核心的利用率
	var err error
	state.PerCPUPercent, err = cpu.Percent(interval, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %w", err)
	}

	// 计算总利用率
	for _, usage := range state.PerCPUPercent {
		state.TotalCPUPercent += usage
	}
	state.TotalCPUPercent /= float64(cores)

	// 获取内存利用率
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}
	state.TotalMemPercent = float64(memInfo.Used) / float64(memInfo.Total) * 100

	// 获取 GPU 利用率
	output, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute nvidia-smi: %w", err)
	}
	state.TotalGPUPercent = parseGPUUtilization(string(output))

	return state, nil
}

func parseGPUUtilization(output string) float64 {
	utilizationStr := strings.TrimSpace(output)
	utilization, err := strconv.ParseFloat(utilizationStr, 64)
	if err != nil {
		fmt.Println("Failed to parse GPU utilization:", err)
		return -1
	}
	return utilization
}

type OSType int

const (
	OSTypeWindows OSType = iota
	OSTypeLinux
)

// GetOSType 返回运行此程序的操作系统类型，如果是非windows，linux操作系统则会返回-1
func GetOSType() OSType {
	os := runtime.GOOS
	switch os {
	case "windows":
		fmt.Println("The operating system is Windows")
		return OSTypeWindows
	case "linux":
		fmt.Println("The operating system is Linux")
		return OSTypeLinux
	default:
		fmt.Printf("The operating system is %s\n", os)
		return -1
	}
}
