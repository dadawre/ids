package ruleStructure

import (
	"fmt"
	"github.com/google/gopacket"
)

// Field 协议字段结构体
type Field struct {
	NameId   int8
	ValuePtr interface{}
}

// Protocol 协议类型及其字段集合结构体
type Protocol struct {
	NameId gopacket.LayerType
	Fields []*Field
}

// Rule 规则结构体
type Rule struct {
	ID        int
	Protocols []*Protocol
	Operator  []int8
	Priority  int
	Action    int8
	Msg       string
	Strings   []*string
}

// RuleNode 代表链表中的一个节点,值为一个规则指针
type RuleNode struct {
	RuleValue *Rule
	Next      *RuleNode
}

// RuleLinkedList 代表链表,指向链表头节点
type RuleLinkedList struct {
	head *RuleNode
}

type RuleOperateError string

func (e RuleOperateError) Error() string {
	return string(e)
}

func (list *RuleLinkedList) GetHead() *RuleNode {
	return list.head
}

// Insert 插入一个新节点到链表中，并保持链表有序
func (list *RuleLinkedList) Insert(value *Rule) {
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
