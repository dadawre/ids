package Action

import (
	"fmt"
	"graduationDesign/back_end/Initialize"
	"os/exec"
)

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
	// 通过 exec.Command 执行 Shell 脚本,向ip_ban_network集合中插入封禁IP网段
	cmd := exec.Command("/bin/bash", Initialize.AddIpNetworkToBanShPath, IPNetworkAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}

func DropIP(IPAddr string) error {
	// 通过 exec.Command 执行 Shell 脚本,向ip_ban集合中插入封禁IP
	cmd := exec.Command("/bin/bash", Initialize.AddIpToBanShPath, IPAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}
