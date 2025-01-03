package Handler

import (
	"fmt"
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

// DropIPv4Packets uses iptables to drop packets based on the provided IPv4 address, ports, and protocol
func DropIPv4Packets(srcIP, dstIP, srcPort, dstPort, protocol string) error {
	cmd := "sudo iptables -A INPUT"

	// Add source IP match if provided
	if srcIP != "" {
		cmd += fmt.Sprintf(" -s %s", srcIP)
	}

	// Add destination IP match if provided
	if dstIP != "" {
		cmd += fmt.Sprintf(" -d %s", dstIP)
	}

	// Add protocol match if provided
	if protocol != "" {
		cmd += fmt.Sprintf(" -p %s", protocol)
	}

	// Add source port match if provided
	if srcPort != "" {
		cmd += fmt.Sprintf(" --sport %s", srcPort)
	}

	// Add destination port match if provided
	if dstPort != "" {
		cmd += fmt.Sprintf(" --dport %s", dstPort)
	}

	// Add the rule to drop the packet
	cmd += " -j DROP"

	// Execute the iptables command
	fmt.Println("Executing command:", cmd) // Print the command for debugging purposes
	return exec.Command("bash", "-c", cmd).Run()
}
