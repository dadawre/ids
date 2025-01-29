package Action

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

// IPType 表示IP地址类型
type IPType int

const (
	Invalid  IPType = iota // 无效格式
	IPv4Addr               // IPv4地址
	IPv6Addr               // IPv6地址
	IPv4Net                // IPv4网段
	IPv6Net                // IPv6网段
)

func CheckIPType(s string) IPType {
	// 尝试解析为CIDR格式
	if ip, ipNet, err := net.ParseCIDR(s); err == nil {
		// 验证CIDR有效性
		if strings.Contains(ip.String(), ".") && !strings.Contains(ip.String(), ":") {
			if ones, _ := ipNet.Mask.Size(); ones > 32 {
				return Invalid
			}
			return IPv4Net
		}
		if strings.Contains(ip.String(), ":") {
			if ones, _ := ipNet.Mask.Size(); ones > 128 {
				return Invalid
			}
			return IPv6Net
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
		return IPv4Addr
	}
	return IPv6Addr
}

// DropMacPackets uses iptables to drop packets based on source MAC, destination MAC, length, and type
func DropMacPackets(srcMAC, dstMAC, length, packetType string) error {
	cmd := "sudo iptables -A INPUT"

	// Add source MAC match if provided
	if srcMAC != "" {
		cmd += fmt.Sprintf(" -m mac --mac-source %s", srcMAC)
	}

	// Add destination MAC match if provided
	if dstMAC != "" {
		cmd += fmt.Sprintf(" -m mac --mac-destination %s", dstMAC)
	}

	// Add packet length match if provided
	if length != "" {
		cmd += fmt.Sprintf(" -m length --length %s", length)
	}

	// Add packet type match if provided
	if packetType != "" {
		cmd += fmt.Sprintf(" -p %s", packetType)
	}

	// Add the rule to drop the packet
	cmd += " -j DROP"

	// Execute the iptables command
	fmt.Println("Executing command:", cmd) // Print the command for debugging purposes
	return exec.Command("bash", "-c", cmd).Run()
}

func DropIPv4Network(IPNetworkAddr string) error {
	//// 通过 exec.Command 执行 Shell 脚本,向ip_ban_network集合中插入封禁IP网段
	//cmd := exec.Command("/bin/bash", Parameter.AddIpNetworkToBanShPath, IPNetworkAddr)
	//output, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return err
	//}
	//// 打印脚本的输出
	//fmt.Println("Output:", string(output))
	return nil
}

func DropIP(IPAddr string) error {
	//// 通过 exec.Command 执行 Shell 脚本,向ip_ban集合中插入封禁IP
	//cmd := exec.Command("/bin/bash", Parameter.AddIpToBanShPath, IPAddr)
	//output, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return err
	//}
	//
	//// 打印脚本的输出
	//fmt.Println("Output:", string(output))
	return nil
}

func UnbanIPv4(IPAddr string) error {
	//cmd := exec.Command("/bin/bash", Parameter.DeleteIpFromShPath, IPAddr)
	//output, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return err
	//}
	//
	//// 打印脚本的输出
	//fmt.Println("Output:", string(output))
	return nil
}

func UnbanIPv4Network(IPAddr string) error {
	//cmd := exec.Command("/bin/bash", Parameter.DeleteIpNetworkFromShPath, IPAddr)
	//output, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return err
	//}
	//
	//// 打印脚本的输出
	//fmt.Println("Output:", string(output))
	return nil
}
