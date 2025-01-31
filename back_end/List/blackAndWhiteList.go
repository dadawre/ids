package List

import (
	"bufio"
	"fmt"
	"graduationDesign/back_end/Action"
	"net"
	"os"
	"strings"
)

type AccessList struct {
	Ips map[string]bool // 存储IP
}

type ListType int

const (
	BlackListType ListType = iota
	WhiteListType
)

// LoadAccessList 加载访问列表文件
func LoadAccessList(filename string, insertType ListType) (*AccessList, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	al := &AccessList{
		Ips: make(map[string]bool),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 尝试解析CIDR格式
		if _, cidr, err := net.ParseCIDR(line); err == nil {
			_ = InsertIPset(cidr.String(), 1, insertType)
			al.Ips[cidr.String()] = true
			continue
		}

		// 尝试解析普通IP格式
		if ip := net.ParseIP(line); ip != nil {
			_ = InsertIPset(ip.String(), 0, insertType)
			al.Ips[ip.String()] = true
			continue
		}

		// 忽略无法解析的行
		fmt.Printf("Warning: 无效条目 %q 已被忽略\n", line)
	}

	return al, scanner.Err()
}

func InsertIPset(insertIP string, IPType int, insertType ListType) error {
	if insertType == 0 { //黑名单
		if IPType == 0 {
			err := Action.DropIP(insertIP)
			if err != nil {
				return err
			}
		} else if IPType == 1 {
			err := Action.DropIPv4Network(insertIP)
			if err != nil {
				return err
			}
		}
	} else if insertType == 1 { //白名单
		if IPType == 0 {
			err := Action.AddWhiteIPv4(insertIP)
			if err != nil {
				return err
			}
		} else if IPType == 1 {
			err := Action.AddWhiteIpv4Network(insertIP)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
