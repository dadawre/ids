package Initialize

import "fmt"

// BanIPDescribe 数据结构
type BanIPDescribe struct {
	BanIP     string `json:"banIP"`
	BanRuleID int    `json:"banRuleID"`
	Describe  string `json:"describe"`
}

// BanIPSet 使用哈希表存储 BanIP 描述信息
type BanIPSet struct {
	data map[string]BanIPDescribe
}

// NewBanIPSet 创建新的 BanIPSet
func NewBanIPSet() *BanIPSet {
	return &BanIPSet{data: make(map[string]BanIPDescribe)}
}

// Insert 插入 BanIPDescribe
func (set *BanIPSet) Insert(banIPDesc BanIPDescribe) bool {
	if _, exists := set.data[banIPDesc.BanIP]; exists {
		return false
	}
	set.data[banIPDesc.BanIP] = banIPDesc
	return true
}

// GetAll 获取所有 BanIPDescribe
func (set *BanIPSet) GetAll() []BanIPDescribe {
	var result []BanIPDescribe
	for _, desc := range set.data {
		result = append(result, desc)
	}
	return result
}

// Delete 删除 BanIPDescribe，根据 BanIP 删除
func (set *BanIPSet) Delete(banIP string) bool {
	if _, exists := set.data[banIP]; exists {
		delete(set.data, banIP)
		return true
	}
	return false
}

// Exists 判断 BanIP 是否已经存在
func (set *BanIPSet) Exists(banIP string) bool {
	_, exists := set.data[banIP]
	return exists
}

// ShowAll 显示所有的 BanIPDescribe
func (set *BanIPSet) ShowAll() {
	for _, desc := range set.data {
		fmt.Printf("BanIP: %s, BanRuleID: %d, Describe: %s\n", desc.BanIP, desc.BanRuleID, desc.Describe)
	}
}
