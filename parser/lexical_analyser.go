package parser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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
	TokenKeywordAny                      // 关键字 any
	TokenKeywordInclude                  // 关键字 include
	TokenKeywordIPv4Var                  // 关键字 ipv4var
	TokenKeywordIPv6Var                  // 关键字 ipv6var
	TokenKeywordPortVar                  // 关键字 portvar
	TokenKeywordVar                      // 关键字 var
	TokenStringValue
	TokenIpv4Value
	TokenIpv4NetworkValue
	TokenIpv6Value
	TokenIpv6NetworkValue
	TokenVarUse
)

const (
	// IPv4地址的正则表达式
	ipv4Pattern = `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

	// IPv4地址段的正则表达式
	ipv4RangePattern = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/(3[0-2]|[12]?[0-9])$`
	ipv6Pattern      = `^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4})|(([0-9a-fA-F]{1,4}:){1,6}:([0-9a-fA-F]{1,4}:){1,6})|(([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,6})|(([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,7})|(([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,8})|(([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,9})|(([0-9a-fA-F]{1,4}:){1}(:[0-9a-fA-F]{1,4}){1,10})|(:((:[0-9a-fA-F]{1,4}){1,11}))$`
)

// ConfigLexer 配置文件的词法分析器
type ConfigLexer struct {
	input      string        // 输入的源代码字符串
	pos        int           // 当前位置
	currLine   int           // 当前行
	curr       rune          // 当前字符
	reader     *bufio.Reader // 用于从标准输入读取输入字符串
	file       *os.File      // File句柄
	CurrTokens []Token       // 存储词法单元的切片
}

// NextToken 获取下一个词法单元
func (l *ConfigLexer) NextToken() Token {
	if l.pos != 0 && l.curr == eof {
		return Token{}
	}
	if l.pos > 0 && !unicode.IsSpace(l.curr) && l.curr != '"' && l.curr != '{' && l.curr != '}' {
		currToken := l.parseIdentifier()
		l.CurrTokens = append(l.CurrTokens, currToken)
		return currToken
	}
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
		if l.curr == '"' || l.curr == '=' || l.curr == '{' || l.curr == '}' {
			continue
		}
		//if unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '$' || l.curr == '/' {
		//
		//}
		if l.curr == eof {
			defer l.file.Close()
			return Token{Type: TokenEOF, Value: ""}
		}

		if unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '$' || l.curr == '/' || l.curr == '_' || l.curr == '.' {
			fmt.Printf("isLetter or isDigit %c\n", l.curr)
			currToken := l.parseIdentifier()
			l.CurrTokens = append(l.CurrTokens, currToken)
			return currToken
		}

		panic(fmt.Sprintf("无效的字符: %c", l.curr))
	}
}

func (l *ConfigLexer) isSpecial(ch rune) bool {
	specials := []rune{'{', '}', '$', '/'}
	for _, s := range specials {
		if ch == s {
			return true
		}
	}
	return false
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
	//fmt.Printf("letter:%c  pos:%d  char:%c\n", l.curr, l.pos, ch)
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
	havePoint := false
	haveBias := false
	haveColon := false

	// 如果首字符为$则代表这是一个标识符
	var haveDollar = l.curr == '$'
	if haveDollar {
		l.readChar()
		for unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '_' {
			l.readChar()
		}
		identifier := l.input[startPos+1 : l.pos-1]
		return Token{Type: TokenVarUse, Value: identifier}
	}

	for unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '.' || l.curr == '/' || l.curr == ':' || l.curr == '_' {
		havePoint = havePoint || l.curr == '.'
		haveBias = haveBias || l.curr == '/'
		haveColon = haveColon || l.curr == ':'
		l.readChar()
	}

	identifier := l.input[startPos : l.pos-1]

	//如果有冒号则代表这是一个Ipv6值
	//如果有冒号则代表这是一个Ipv6值

	// 匹配IPv6地址
	reIPv6 := regexp.MustCompile(ipv6Pattern)

	if reIPv6.MatchString(identifier) {
		fmt.Println(identifier, "   ", haveBias)
		if !haveBias {
			return Token{Type: TokenIpv6Value, Value: identifier}
		} else {
			return Token{Type: TokenIpv6NetworkValue, Value: identifier}
		}
	}

	reIPv4Net := regexp.MustCompile(ipv4RangePattern)
	if reIPv4Net.MatchString(identifier) {
		fmt.Println(identifier, "   ", haveBias)
		return Token{Type: TokenIpv4NetworkValue, Value: identifier}
	}

	// 匹配IPv4地址
	reIPv4 := regexp.MustCompile(ipv4Pattern)
	if reIPv4.MatchString(identifier) {
		fmt.Println(identifier, "   ", haveBias)
		return Token{Type: TokenIpv4Value, Value: identifier}
	}

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

		CurrTokensLen := len(l.CurrTokens)
		if CurrTokensLen != 0 && l.CurrTokens[CurrTokensLen-1].Type == TokenIdentifier {
			return Token{Type: TokenStringValue, Value: identifier}
		}
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
	initializeTokens := make([]Token, 0, 0)
	return &ConfigLexer{reader: reader, input: str, CurrTokens: initializeTokens}
}

// GetAllToken 对输入的源代码进行词法分析，返回词法单元序列
func (l *ConfigLexer) GetAllToken() []Token {
	for {
		myToken := l.NextToken()
		if myToken.Type == TokenEOF {
			break
		}
	}
	return l.CurrTokens
}
