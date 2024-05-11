package parser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"
)

// Token 表示词法单元
type Token struct {
	Type  TokenType // 单词类型
	Value string    // 单词的值
}

// TokenType 表示单词类型
type TokenType int

// 定义词法单元类型的常量
const (
	TokenEOF            TokenType = iota // 文件结束
	TokenIdentifier                      // 标识符
	TokenKeywordConfig                   // 关键字 config
	TokenKeywordInclude                  // 关键字 include
	TokenKeywordAny                      // 关键字 any
	TokenKeywordIPv4Var                  // 关键字 ipv4var
	TokenKeywordIPv6Var                  // 关键字 ipv6var
	TokenKeywordPortVar                  // 关键字 portvar
	TokenKeywordVar                      // 关键字 var
	TokenComment                         // 注释
)

// ConfigLexer 配置文件的词法分析器
type ConfigLexer struct {
	input    string        // 输入的源代码字符串
	pos      int           // 当前位置
	currLine int           // 当前行
	curr     rune          // 当前字符
	reader   *bufio.Reader // 用于从标准输入读取输入字符串
	file     *os.File
	tokens   []Token // 存储词法单元的切片
}

// NextToken 获取下一个词法单元
func (l *ConfigLexer) NextToken() Token {
	for {
		l.readChar()
		if unicode.IsSpace(l.curr) {
			fmt.Println("isWhiteSpace")
			l.skipWhitespace()
		}
		if l.curr == '#' {
			fmt.Println("isComment")
			l.skipComment()
			continue
		}
		if unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) {
			fmt.Printf("isLetter or isDigit %c\n", l.curr)
			return l.parseIdentifier()
		}
		if l.curr == '"' {
			continue
		}
		if l.curr == eof {
			defer l.file.Close()
			return Token{Type: TokenEOF, Value: ""}
		}
		panic(fmt.Sprintf("无效的字符: %c", l.curr))
	}
}

// 读取下一个字符
func (l *ConfigLexer) readChar() {
	//读取一个字符
	ch, _, err := l.reader.ReadRune()

	if err != nil {
		l.curr = eof
		fmt.Println(err)
		return
	}
	l.curr = ch
	l.pos++
	//fmt.Println("letter:", l.curr, "  pos:", l.pos, "  char:", ch)
	fmt.Printf("letter:%c  pos:%d  char:%c\n", l.curr, l.pos, ch)
}

// 跳过空白字符
func (l *ConfigLexer) skipWhitespace() {
	for unicode.IsSpace(l.curr) {
		l.readChar()
	}
}

// 跳过注释
func (l *ConfigLexer) skipComment() {
	for l.curr != '\n' && l.curr != eof {
		l.readChar()
	}
}

// 解析标识符或关键字
func (l *ConfigLexer) parseIdentifier() Token {
	startPos := l.pos - 1
	for unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) {
		l.readChar()
	}
	identifier := l.input[startPos : l.pos-1]
	switch strings.ToLower(identifier) {
	case "config":
		return Token{Type: TokenKeywordConfig, Value: identifier}
	case "include":
		return Token{Type: TokenKeywordInclude, Value: identifier}
	case "any":
		return Token{Type: TokenKeywordAny, Value: identifier}
	case "ipv4var":
		return Token{Type: TokenKeywordIPv4Var, Value: identifier}
	case "ipv6var":
		return Token{Type: TokenKeywordIPv6Var, Value: identifier}
	case "portvar":
		return Token{Type: TokenKeywordPortVar, Value: identifier}
	case "var":
		return Token{Type: TokenKeywordVar, Value: identifier}
	default:
		return Token{Type: TokenIdentifier, Value: identifier}
	}
}

const eof = rune(0)

// LexerInterface 语法分析接口
type LexerInterface interface {
	Tokenize() []Token
}

// NewConfigLexer 创建一个新的词法分析器
func NewConfigLexer(path string) *ConfigLexer {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	// 创建bufio.Reader之前检查文件打开是否成功
	info, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return nil
	}
	fmt.Println(info)
	// 文件打开成功，创建bufio.Reader
	reader := bufio.NewReader(file)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading file:", err)
		return nil
	}

	// 将字节切片转换成字符串
	str := string(content)

	return &ConfigLexer{reader: reader, input: str}
}

// Tokenize 对输入的源代码进行词法分析，返回词法单元序列
func (l *ConfigLexer) Tokenize() []Token {
	l.tokens = make([]Token, 10)

	return l.tokens
}
