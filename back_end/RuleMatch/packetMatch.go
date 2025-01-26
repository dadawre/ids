package RuleMatch

import (
	"errors"
)

// EvaluateExpression 评估基于布尔值数组和操作符数组的逻辑表达式。
// 参数：
// bools - 布尔值数组，例如 [true, false, true]。
// ops - 操作符数组，其中：
//
//	1 表示 && (与操作符)
//	2 表示 || (或操作符)
//	3 表示 ( (左括号)
//	4 表示 ) (右括号)
//
// 返回：
// 布尔值表示表达式的评估结果，如果表达式无效则返回错误。
func EvaluateExpression(bools []bool, ops []int8) (bool, error) {
	var stack []bool   // 用于存储布尔值的栈
	var opStack []int8 // 用于存储操作符的栈

	// 辅助函数：应用栈顶的操作符到栈顶的两个布尔值
	applyOperator := func() error {
		if len(stack) < 2 {
			return errors.New("表达式无效：操作数不足")
		}
		b2 := stack[len(stack)-1]    // 栈顶的第二个布尔值
		b1 := stack[len(stack)-2]    // 栈顶的第一个布尔值
		stack = stack[:len(stack)-2] // 弹出两个布尔值

		op := opStack[len(opStack)-1]      // 获取栈顶的操作符
		opStack = opStack[:len(opStack)-1] // 弹出操作符

		var result bool
		if op == 1 { // && 操作符
			result = b1 && b2
		} else if op == 2 { // || 操作符
			result = b1 || b2
		} else {
			return errors.New("无效的操作符")
		}

		stack = append(stack, result) // 将结果压入栈中
		return nil
	}

	boolIndex := 0 // 当前布尔值数组的索引
	for _, op := range ops {
		if op == 3 { // 左括号 (
			opStack = append(opStack, op) // 左括号直接压入操作符栈
		} else if op == 4 { // 右括号 )
			// 一直弹出操作符并应用，直到遇到左括号
			for len(opStack) > 0 && opStack[len(opStack)-1] != 3 {
				if err := applyOperator(); err != nil {
					return false, err
				}
			}
			if len(opStack) == 0 || opStack[len(opStack)-1] != 3 {
				return false, errors.New("括号不匹配")
			}
			opStack = opStack[:len(opStack)-1] // 弹出左括号
		} else if op == 1 || op == 2 { // && 或 ||
			// 处理当前操作符的优先级，弹出并应用栈顶优先级相等或更高的操作符
			for len(opStack) > 0 && opStack[len(opStack)-1] <= 2 {
				if err := applyOperator(); err != nil {
					return false, err
				}
			}
			opStack = append(opStack, op) // 当前操作符压入栈
		} else if op == 0 { // 操作数占位符
			if boolIndex >= len(bools) {
				return false, errors.New("布尔值不足")
			}
			stack = append(stack, bools[boolIndex]) // 将布尔值压入栈
			boolIndex++
		} else {
			return false, errors.New("无效的操作符或操作数")
		}
	}

	// 处理所有剩余的操作符
	for len(opStack) > 0 {
		if opStack[len(opStack)-1] == 3 {
			return false, errors.New("括号不匹配")
		}
		if err := applyOperator(); err != nil {
			return false, err
		}
	}

	// 检查栈中是否仅剩一个结果
	if len(stack) != 1 {
		return false, errors.New("表达式无效")
	}

	return stack[0], nil // 返回最终的评估结果
}
