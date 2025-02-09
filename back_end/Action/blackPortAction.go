package Action

import (
	"fmt"
	"graduationDesign/back_end/Initialize/Parameter"
	"os/exec"
)

func AddBlackPort(port string) error {
	// 通过 exec.Command 执行 Shell 脚本,向port_ban集合中插入封禁port
	cmd := exec.Command("/bin/bash", Parameter.AddPortToBanShPath, port)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}

func DeleteBlackPort(port string) error {
	// 通过 exec.Command 执行 Shell 脚本,向port_ban集合中插入封禁port
	cmd := exec.Command("/bin/bash", Parameter.DeleteBanPortFromBanSh, port)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 打印脚本的输出
	fmt.Println("Output:", string(output))
	return nil
}
