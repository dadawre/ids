package Action

import (
	"fmt"
	"graduationDesign/back_end/Initialize/Parameter"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

// IPType 表示IP地址类型
type IPType int

const (
	Invalid      IPType = iota // 无效格式
	IPv4AddrType               // IPv4地址
	IPv6AddrType               // IPv6地址
	IPv4NetType                // IPv4网段
	IPv6NetType                // IPv6网段
	PortType
)

// CheckIPType 查看字符串类型
func CheckIPType(s string) IPType {
	if _, err := strconv.Atoi(s); err == nil {
		return PortType
	}

	// 尝试解析为CIDR格式
	if ip, ipNet, err := net.ParseCIDR(s); err == nil {
		// 验证CIDR有效性
		if strings.Contains(ip.String(), ".") && !strings.Contains(ip.String(), ":") {
			if ones, _ := ipNet.Mask.Size(); ones > 32 {
				return Invalid
			}
			return IPv4NetType
		}
		if strings.Contains(ip.String(), ":") {
			if ones, _ := ipNet.Mask.Size(); ones > 128 {
				return Invalid
			}
			return IPv6NetType
		}
		return Invalid
	}

	// 尝试解析为普通IP地址
	ip := net.ParseIP(s)
	if ip == nil {
		return Invalid
	}

	// 判断IP版本
	if ip.To4() != nil {
		return IPv4AddrType
	}
	return IPv6AddrType
}

func DropIPv4Network(IPNetworkAddr ...string) error {
	// 通过 exec.Command 执行 Shell 脚本,向ip_ban_network集合中插入封禁IP网段
	var cmd *exec.Cmd
	if len(IPNetworkAddr) == 1 {
		cmd = exec.Command("/bin/bash", Parameter.AddIpNetworkToBanShPath, IPNetworkAddr[0])
	} else if len(IPNetworkAddr) >= 2 {
		cmd = exec.Command("/bin/bash", Parameter.AddIpNetworkToBanShPath, IPNetworkAddr[0], IPNetworkAddr[1])
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}

func DropIP(IPAddr ...string) error {
	// 通过 exec.Command 执行 Shell 脚本,向ip_ban集合中插入封禁IP
	var cmd *exec.Cmd
	if len(IPAddr) == 1 {
		cmd = exec.Command("/bin/bash", Parameter.AddIpNetworkToBanShPath, IPAddr[0])
	} else if len(IPAddr) >= 2 {
		cmd = exec.Command("/bin/bash", Parameter.AddIpNetworkToBanShPath, IPAddr[0], IPAddr[1])
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}

func UnbanIPv4(IPAddr string) error {
	cmd := exec.Command("/bin/bash", Parameter.DeleteIpFromShPath, IPAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}

func UnbanIPv4Network(IPAddr string) error {
	cmd := exec.Command("/bin/bash", Parameter.DeleteIpNetworkFromShPath, IPAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}
