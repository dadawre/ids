package parser

// Token 表示词法单元
type Token struct {
	Type  TokenType // 单词类型
	Value string    // 单词的值
}

// TokenType 表示单词类型
type TokenType int

// 定义一些单词类型的常量
const (
	TokenTypeIdentifier TokenType = iota // 标识符
	TokenTypeKeyword                     // 关键字
	TokenTypeOperator                    // 操作符
	// 其他类型...
)

// ConfigLexer 表示词法分析器
type ConfigLexer struct {
	input  string  // 输入的源代码字符串
	tokens []Token // 存储词法单元的切片
}

type LexerInterface interface {
}

// NewLexer 创建一个新的词法分析器
func NewLexer(input string) *ConfigLexer {
	return &ConfigLexer{input: input}
}

// Tokenize 对输入的源代码进行词法分析，返回词法单元序列
func (l *ConfigLexer) Tokenize() []Token {
	// 在这里进行词法分析，生成词法单元序列，并存储到 l.tokens 中
	// ...
	return l.tokens
}
