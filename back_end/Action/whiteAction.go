package Action

import (
	"fmt"
	"graduationDesign/back_end/Initialize/Parameter"
	"os/exec"
)

func AddWhiteIpv4Network(IPNetworkAddr string) error {
	// 通过 exec.Command 执行 Shell 脚本,向ip_ban_network集合中插入封禁IP网段
	cmd := exec.Command("/bin/bash", Parameter.AddIpToWhiteListNetworkShPath, IPNetworkAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}

func AddWhiteIPv4(IPAddr string) error {
	// 通过 exec.Command 执行 Shell 脚本,向ip_ban集合中插入封禁IP
	cmd := exec.Command("/bin/bash", Parameter.AddIpToWhiteListShPath, IPAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}

func DeleteWhiteIPv4Network(IPNetworkAddr string) error {
	cmd := exec.Command("/bin/bash", Parameter.DeleteWhiteIpNetworkSh, IPNetworkAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}

func DeleteWhiteIPv4(IPAddr string) error {
	cmd := exec.Command("/bin/bash", Parameter.DeleteWhiteIpSh, IPAddr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}
