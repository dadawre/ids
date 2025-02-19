package Structure

import (
	"fmt"
	"graduationDesign/back_end/Action"
	"log"
	"strconv"
	"sync"
	"time"
)

// BanIPDescribe 数据结构
type BanIPDescribe struct {
	BanIP     string    `json:"banIP"`
	BanRuleID int       `json:"banRuleID"`
	Describe  string    `json:"describe"`
	UnbanTime time.Time `json:"unbanTime"` // 解封时间字段
	BanPeriod int       `json:"banPeriod"` //封禁时长
}

// BanIPSet 使用哈希表存储 BanIP 描述信息
type BanIPSet struct {
	Data        map[string]BanIPDescribe
	timerGroups map[time.Time]*timerGroup // 管理定时器组
	mu          sync.Mutex                //锁
	allowedDiff time.Duration             // 允许的时间误差
}

// 定时器结构体
type timerGroup struct {
	ips   []string
	timer *time.Timer
}

// NewBanIPSet 创建新的 BanIPSet，需传入允许的时间误差
func NewBanIPSet(allowedDiff time.Duration) *BanIPSet {
	return &BanIPSet{
		Data:        make(map[string]BanIPDescribe),
		timerGroups: make(map[time.Time]*timerGroup),
		allowedDiff: allowedDiff,
	}
}

// Insert 插入 BanIPDescribe，自动处理时间合并逻辑
func (set *BanIPSet) Insert(banIPDesc BanIPDescribe) bool {
	set.mu.Lock()
	defer set.mu.Unlock()

	if _, exists := set.Data[banIPDesc.BanIP]; exists {
		return false
	}

	// 查找可合并的时间组
	var foundTime time.Time
	for t := range set.timerGroups {
		if absDuration(t.Sub(banIPDesc.UnbanTime)) <= set.allowedDiff {
			foundTime = t
			break
		}
	}

	if !foundTime.IsZero() {
		// 合并到现有组
		banIPDesc.UnbanTime = foundTime
		set.timerGroups[foundTime].ips = append(set.timerGroups[foundTime].ips, banIPDesc.BanIP)
	} else {
		// 创建新时间组
		foundTime = banIPDesc.UnbanTime
		group := &timerGroup{ips: []string{banIPDesc.BanIP}}
		set.timerGroups[foundTime] = group

		durationUntil := time.Until(foundTime)
		if durationUntil <= 0 {
			set.unbanIPs(foundTime)
		} else {
			group.timer = time.AfterFunc(durationUntil, func() {
				set.unbanIPs(foundTime)
			})
		}
	}

	set.Data[banIPDesc.BanIP] = banIPDesc

	//从IPset中插入
	IPType := Action.CheckIPType(banIPDesc.BanIP)
	var err error
	if IPType == Action.IPv4AddrType {
		err = Action.DropIP(banIPDesc.BanIP, strconv.Itoa(banIPDesc.BanPeriod*60))
	} else if IPType == Action.IPv4NetType {
		err = Action.DropIPv4Network(banIPDesc.BanIP, strconv.Itoa(banIPDesc.BanPeriod*60))
	}

	if err != nil {
		log.Fatal(err)
	}

	return true
}

// Delete 删除封禁IP并清理定时器
func (set *BanIPSet) Delete(banIP string) bool {
	set.mu.Lock()
	defer set.mu.Unlock()

	desc, exists := set.Data[banIP]
	if !exists {
		return false
	}

	delete(set.Data, banIP)

	// 从对应时间组中移除
	if group, ok := set.timerGroups[desc.UnbanTime]; ok {

		//从IPset中删除
		IPType := Action.CheckIPType(banIP)
		var err error
		if IPType == Action.IPv4AddrType {
			err = Action.UnbanIPv4(banIP)
		} else if IPType == Action.IPv4NetType {
			err = Action.UnbanIPv4Network(banIP)
		}

		if err != nil {
			log.Fatal(err)
		}

		//从定时器组中删除
		for i, ip := range group.ips {
			if ip == banIP {
				group.ips = append(group.ips[:i], group.ips[i+1:]...)
				break
			}
		}

		// 清理空组
		if len(group.ips) == 0 {
			if group.timer != nil {
				group.timer.Stop()
			}
			delete(set.timerGroups, desc.UnbanTime)
		}
	}

	return true
}

// 解封指定时间组的所有IP
func (set *BanIPSet) unbanIPs(baseTime time.Time) {
	set.mu.Lock()
	defer set.mu.Unlock()

	if group, exists := set.timerGroups[baseTime]; exists {
		for _, ip := range group.ips {
			IPType := Action.CheckIPType(ip)
			var err error
			if IPType == Action.IPv4AddrType {
				err = Action.UnbanIPv4(ip)
			} else if IPType == Action.IPv4NetType {
				err = Action.UnbanIPv4Network(ip)
			}
			delete(set.Data, ip)
			if err != nil {
				return
			}
		}
		delete(set.timerGroups, baseTime)
	}
}

// 辅助函数：获取时间间隔绝对值
func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

func (set *BanIPSet) GetAllTransformArray() []BanIPDescribe {
	var result []BanIPDescribe
	for _, desc := range set.Data {
		result = append(result, desc)
	}
	return result
}

func (set *BanIPSet) Exists(banIP string) bool {
	_, exists := set.Data[banIP]
	return exists
}

func (set *BanIPSet) ShowAll() {
	for _, desc := range set.Data {
		fmt.Printf("BanIP: %s, BanRuleID: %d, Describe: %s, UnbanTime: %s\n",
			desc.BanIP, desc.BanRuleID, desc.Describe, desc.UnbanTime.Format(time.RFC3339))
	}
}
