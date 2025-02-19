package List

import (
	"bufio"
	"fmt"
	"graduationDesign/back_end/Action"
	"net"
	"os"
	"strconv"
	"strings"
)

type AccessList struct {
	IDindex int                         // Used to assign a unique ID; only increases.
	Addrs   map[string]WBValueStructure // Stores entries using the address as key.
}

type WBValueStructure struct {
	ID        int        `json:"id"`
	Addr      string     `json:"addr"`      // Address string
	IsWhite   bool       `json:"isWhite"`   // Whether it is white-listed
	ValueType WBAddrType `json:"valueType"` // Address type
}

type WBAddrType int

const (
	BlackIPv4Type WBAddrType = iota
	WhiteIPv4Type
	BlackIPv4NetworkType
	WhiteIPv4NetworkType
	BlackPortType
)

// LoadAccessList 加载访问列表文件
func LoadAccessList(filename string, isWhite bool, AppendList *AccessList) (*AccessList, error) {
	fmt.Printf(filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var al *AccessList
	if AppendList == nil {
		al = &AccessList{
			Addrs:   make(map[string]WBValueStructure),
			IDindex: 0,
		}
	} else {
		al = AppendList
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var WBValue WBValueStructure
		WBValue.IsWhite = isWhite

		// 尝试解析CIDR格式
		if _, cidr, err := net.ParseCIDR(line); err == nil {
			if isWhite == true {
				_ = InsertIPset(cidr.String(), WhiteIPv4NetworkType)
				WBValue.ValueType = WhiteIPv4NetworkType
			} else {
				_ = InsertIPset(cidr.String(), BlackIPv4NetworkType)
				WBValue.ValueType = BlackIPv4NetworkType
			}
			WBValue.Addr = cidr.String()
			WBValue.IsWhite = isWhite
			WBValue.ID = al.IDindex
			al.Addrs[cidr.String()] = WBValue
			al.IDindex++
			continue
		}

		// 尝试解析普通IP格式
		if ip := net.ParseIP(line); ip != nil {
			if isWhite == true {
				_ = InsertIPset(ip.String(), WhiteIPv4Type)
				WBValue.ValueType = WhiteIPv4Type
			} else {
				_ = InsertIPset(ip.String(), BlackIPv4Type)
				WBValue.ValueType = BlackIPv4Type
			}
			WBValue.Addr = ip.String()
			WBValue.IsWhite = isWhite
			WBValue.ID = al.IDindex
			al.Addrs[ip.String()] = WBValue
			al.IDindex++
			continue
		}

		//解析port格式
		if port, err := strconv.Atoi(line); err == nil && port > 0 && port < 65536 {
			portStr := strconv.Itoa(port)
			WBValue.Addr = portStr
			WBValue.IsWhite = isWhite
			WBValue.ValueType = BlackPortType
			WBValue.ID = al.IDindex
			al.Addrs[portStr] = WBValue
			al.IDindex++
			_ = InsertIPset(portStr, BlackPortType)
			continue
		}

		// 忽略无法解析的行
		fmt.Printf("Warning: 无效条目 %q 已被忽略\n", line)
	}

	return al, scanner.Err()
}

func InsertIPset(insertAddr string, insertType WBAddrType) error {
	var err error = nil
	if insertType == BlackIPv4Type { //黑名单
		err = Action.DropIP(insertAddr)
	} else if insertType == WhiteIPv4Type { //白名单
		err = Action.AddWhiteIPv4(insertAddr)
	} else if insertType == BlackPortType {
		err = Action.AddBlackPort(insertAddr)
	} else if insertType == BlackIPv4NetworkType {
		err = Action.DropIPv4Network(insertAddr)
	} else if insertType == WhiteIPv4NetworkType {
		err = Action.AddWhiteIpv4Network(insertAddr)
	}

	return err
}

func DeleteIPset(insertAddr string, insertType WBAddrType) error {
	var err error = nil
	if insertType == BlackIPv4Type { //黑名单
		err = Action.UnbanIPv4(insertAddr)
	} else if insertType == WhiteIPv4Type { //白名单
		err = Action.DeleteWhiteIPv4(insertAddr)
	} else if insertType == BlackPortType {
		err = Action.DeleteBlackPort(insertAddr)
	} else if insertType == BlackIPv4NetworkType {
		err = Action.UnbanIPv4Network(insertAddr)
	} else if insertType == WhiteIPv4NetworkType {
		err = Action.DeleteWhiteIPv4Network(insertAddr)
	}

	return err
}
