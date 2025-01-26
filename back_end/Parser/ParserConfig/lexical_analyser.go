package ParserConfig

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// Token 代表一个词法标志
type Token struct {
	Type  TokenType // Token type
	Value string    // Token value
}

// TokenType 标志类型
type TokenType int

// 不同标记类型的常量。
const (
	TokenEOF              TokenType = iota //结尾标志
	TokenIdentifier                        //标识符
	TokenKeywordConfig                     //config关键字
	TokenKeywordAny                        //any关键字
	TokenKeywordInclude                    //include关键字
	TokenKeywordIPv4Var                    //IPv4关键字
	TokenKeywordIPv6Var                    //IPv6关键字
	TokenKeywordPortVar                    //port关键字
	TokenKeywordVar                        //var关键字
	TokenStringValue                       //字符串值
	TokenIpv4Value                         //IPv4值
	TokenIpv4NetworkValue                  //IPv4网络块值
	TokenIpv6Value                         //IPv6值
	TokenIpv6NetworkValue                  //IPv6网络块值
	TokenVarUse                            //var值
)

// IP正则表达式
const (
	ipv4Pattern      = `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	ipv4RangePattern = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/(3[0-2]|[12]?[0-9])$`
	ipv6Pattern      = `^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4})|(([0-9a-fA-F]{1,4}:){1,6}:([0-9a-fA-F]{1,4}:){1,6})|(([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,6})|(([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,7})|(([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,8})|(([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,9})|(([0-9a-fA-F]{1,4}:){1}(:[0-9a-fA-F]{1,4}){1,10})|(:((:[0-9a-fA-F]{1,4}){1,11}))$`
)

// ConfigLexer 配置文件词法解析
type ConfigLexer struct {
	input      string //
	pos        int    //当前解析位置
	currLine   int    //当前解析行
	curr       rune
	reader     *bufio.Reader
	file       *os.File
	CurrTokens []Token
}

// 自定义解析错误类型。
type currCharError string

func (e *currCharError) Error() string {
	return fmt.Sprintf("currCharError: %s", string(*e))
}

// NextToken 解析一个标志
func (l *ConfigLexer) NextToken() (Token, error) {
	if l.pos != 0 && l.curr == eof {
		err := currCharError("curr is EOF")
		return Token{}, &err
	}
	if l.pos > 0 && !unicode.IsSpace(l.curr) && l.curr != '"' && l.curr != '{' && l.curr != '}' {
		currToken := l.parseIdentifier()
		l.CurrTokens = append(l.CurrTokens, currToken)
		return currToken, nil
	}
	for {
		l.readChar()
		if unicode.IsSpace(l.curr) {
			l.skipWhitespace()
		}
		if l.curr == '#' {
			l.skipComment()
			continue
		}
		if l.curr == '"' || l.curr == '=' || l.curr == '{' || l.curr == '}' {
			continue
		}
		if l.curr == eof {
			break
		}

		if unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '$' || l.curr == '/' || l.curr == '_' || l.curr == '.' || l.curr == '[' || l.curr == ']' || l.curr == ',' || l.curr == '@' {
			currToken := l.parseIdentifier()
			l.CurrTokens = append(l.CurrTokens, currToken)
			return currToken, nil
		}

		panic(fmt.Sprintf("Invalid character: %c", l.curr))
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(l.file)
	err := currCharError("curr is EOF")
	return Token{Type: TokenEOF, Value: ""}, &err

}

// isSpecial 是否是特殊字符
func (l *ConfigLexer) isSpecial(ch rune) bool {
	specials := []rune{'{', '}', '$'}
	for _, s := range specials {
		if ch == s {
			return true
		}
	}
	return false
}

// readChar 读出下一个字符
func (l *ConfigLexer) readChar() {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		l.curr = eof
		return
	}
	l.curr = ch
	l.pos++
}

// skipWhitespace 跳空格
func (l *ConfigLexer) skipWhitespace() {
	for unicode.IsSpace(l.curr) {
		l.readChar()
	}
}

// skipComment 跳注释
func (l *ConfigLexer) skipComment() {
	for l.curr != '\n' && l.curr != eof {
		//fmt.Println(l.curr)
		l.readChar()
	}
}

// parseIdentifier 解析标识符或关键字。
func (l *ConfigLexer) parseIdentifier() Token {
	startPos := l.pos - 1
	havePoint, haveBias, haveColon := false, false, false

	var haveDollar = l.curr == '$'
	if haveDollar {
		l.readChar()
		for unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '_' || l.curr == '@' {
			l.readChar()
		}
		identifier := l.input[startPos+1 : l.pos-1]
		return Token{Type: TokenVarUse, Value: identifier}
	}

	for unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '.' || l.curr == '/' || l.curr == ':' || l.curr == '_' || l.curr == '[' || l.curr == ']' || l.curr == ',' || l.curr == '@' {
		havePoint = havePoint || l.curr == '.'
		haveBias = haveBias || l.curr == '/'
		haveColon = haveColon || l.curr == ':'
		l.readChar()
	}

	identifier := l.input[startPos : l.pos-1]

	reIPv6 := regexp.MustCompile(ipv6Pattern)
	if reIPv6.MatchString(identifier) {
		if !haveBias {
			return Token{Type: TokenIpv6Value, Value: identifier}
		} else {
			return Token{Type: TokenIpv6NetworkValue, Value: identifier}
		}
	}

	reIPv4Net := regexp.MustCompile(ipv4RangePattern)
	if reIPv4Net.MatchString(identifier) {
		return Token{Type: TokenIpv4NetworkValue, Value: identifier}
	}

	reIPv4 := regexp.MustCompile(ipv4Pattern)
	if reIPv4.MatchString(identifier) {
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
		if CurrTokensLen != 0 && (l.CurrTokens[CurrTokensLen-1].Type == TokenIdentifier || l.CurrTokens[CurrTokensLen-1].Type == TokenKeywordInclude ||
			l.CurrTokens[CurrTokensLen-1].Type == TokenVarUse) {
			return Token{Type: TokenStringValue, Value: identifier}
		}
		return Token{Type: TokenIdentifier, Value: identifier}
	}
}

const eof = rune(0)

//// LexerInterface represents the interface for a lexer.
//type LexerInterface interface {
//	Tokenize() []Token
//}

// NewConfigLexer 为给定的文件路径创建一个新的词法分析器
func NewConfigLexer(path string) *ConfigLexer {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return nil
	}

	reader := bufio.NewReader(file)
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return nil
	}

	str := string(content)
	return &ConfigLexer{reader: reader, input: str, file: file, CurrTokens: []Token{}}
}

// GetAllToken 返回解析出的所有标志
func (l *ConfigLexer) GetAllToken() []Token {
	for {
		myToken, _ := l.NextToken()
		//fmt.Println(myToken.Type, "   ", myToken.Value)
		if myToken.Type == TokenEOF {
			break
		}
	}
	return l.CurrTokens
}
