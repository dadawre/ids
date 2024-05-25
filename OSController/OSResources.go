package OSController

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var currentTotalCPUPercent = 0.0
var currentPerCPUPercent []float64
var currentTotalMemPercent = 0.0
var currentTotalGPUPercent = 0.0

func GetOSStates() {
	interval := 1 * time.Second

	// 获取 CPU 核心数
	cores, err := cpu.Counts(true)
	if err != nil {
		fmt.Println("Failed to get CPU count:", err)
		return
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Failed to get memory info:", err)
		return
	}

	// 创建一个每秒触发一次的定时器
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		currentTotalCPUPercent = 0.0
		currentTotalGPUPercent = 0.0
		currentTotalMemPercent = 0.0

		// 获取每个核心的利用率
		currentPerCPUPercent, err = cpu.Percent(interval, true)
		if err != nil {
			fmt.Println("Failed to get CPU usage:", err)
			currentPerCPUPercent = nil
			currentTotalCPUPercent = -1
		}

		// 计算总利用率
		for _, usage := range currentPerCPUPercent {
			currentTotalCPUPercent += usage
		}
		currentTotalCPUPercent /= float64(cores)

		//// 输出总利用率
		fmt.Printf("Total CPU usage: %.2f%%\n", currentTotalCPUPercent)
		fmt.Println(currentPerCPUPercent)

		// 计算总内存利用率
		currentTotalMemPercent = float64(memInfo.Used) / float64(memInfo.Total) * 100
		// 输出总内存利用率
		fmt.Printf("Total memory usage: %.2f%%\n", currentTotalMemPercent)

		output, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits").Output()
		if err != nil {
			fmt.Println("Failed to execute nvidia-smi:", err)
			return
		}
		// 解析输出并提取 GPU 利用率信息
		gpuUtilization := strings.TrimSpace(string(output))

		// 输出 GPU 利用率
		fmt.Println("GPU Utilization:", gpuUtilization, "%")
	}
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
