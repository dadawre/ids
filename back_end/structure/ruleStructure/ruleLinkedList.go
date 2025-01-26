package ruleStructure

import (
	"fmt"
	"github.com/google/gopacket"
)

// Field 协议字段结构体
type Field struct {
	NameId   int8        `json:"fieldNameId"`
	ValuePtr interface{} `json:"value_ptr"`
}

// Protocol 协议类型及其字段集合结构体
type Protocol struct {
	NameId gopacket.LayerType `json:"protocolNameId"`
	Fields []Field            `json:"fields"`
}

// RuleNodeValue 规则结构体
type RuleNodeValue struct {
	ID            int
	Protocols     []Protocol
	Operators     []int8
	Priority      int
	Action        int8
	Msg           string
	StringIDs     []int
	NetworkValues RuleNodeStringValue
}

// RuleNodeStringValue 游览器显示规则结构体
type RuleNodeStringValue struct {
	ID                     int      `json:"ID"`
	ProtocolsString        string   `json:"ProtocolsString"`
	Priority               int      `json:"Priority"`
	ActionString           string   `json:"ActionString"`
	MsgString              string   `json:"MsgString"`
	StringParameterStrings []string `json:"StringParameterStrings"`
}

// RuleNode 代表链表中的一个节点,值为一个规则指针
type RuleNode struct {
	RuleValue *RuleNodeValue
	Next      *RuleNode
}

// RuleLinkedList 代表链表,指向链表头节点
type RuleLinkedList struct {
	head *RuleNode
}

type RuleOperateError string

var StringArray []RuleNodeStringValue = make([]RuleNodeStringValue, 0)

func (e RuleOperateError) Error() string {
	return string(e)
}

func (list *RuleLinkedList) GetHead() *RuleNode {
	return list.head
}

// Insert 插入一个新节点到链表中，并保持链表有序
func (list *RuleLinkedList) Insert(value *RuleNodeValue) {
	newNode := &RuleNode{RuleValue: value}
	if list.head == nil || list.head.RuleValue.Priority >= value.Priority {
		newNode.Next = list.head
		list.head = newNode
		return
	}
	current := list.head
	for current.Next != nil && current.Next.RuleValue.Priority < value.Priority {
		current = current.Next
	}
	newNode.Next = current.Next
	current.Next = newNode
}

// Delete 删除指定ID
func (list *RuleLinkedList) Delete(ruleID int) error {
	if list.head == nil {
		return RuleOperateError("LinkEmpty")
	}
	if list.head.RuleValue.ID == ruleID {
		list.head = list.head.Next
		return nil
	}
	current := list.head
	for current.Next != nil && current.Next.RuleValue.ID != ruleID {
		current = current.Next
	}
	if current.Next != nil {
		current.Next = current.Next.Next
		return nil
	}
	return RuleOperateError("NoSpecificIdInLink")
}

// Traverse 遍历并打印链表中的所有节点值
func (list *RuleLinkedList) Traverse() {
	current := list.head
	for current != nil {
		fmt.Print(current.RuleValue, " ")
		current = current.Next
	}
	fmt.Println()
}

func (list *RuleLinkedList) GetNodeById(targetID int) *RuleNode {
	current := list.head
	for current != nil {
		if current.RuleValue.ID == targetID {
			return current
		}
		current = current.Next
	}
	return nil
}

func (list *RuleLinkedList) GetStringArray() []RuleNodeStringValue {
	if len(StringArray) == 0 {
		current := list.head
		for current != nil {
			StringArray = append(StringArray, current.RuleValue.NetworkValues)
			//fmt.Println(current.RuleValue.NetworkValues)
			current = current.Next
		}
	}
	fmt.Println(StringArray)
	return StringArray
}
