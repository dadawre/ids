package parserConfig

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// Token represents a lexical token.
type Token struct {
	Type  TokenType // Token type
	Value string    // Token value
}

// TokenType represents the type of a token.
type TokenType int

// Constants for different token types.
const (
	TokenEOF TokenType = iota
	TokenIdentifier
	TokenKeywordConfig
	TokenKeywordAny
	TokenKeywordInclude
	TokenKeywordIPv4Var
	TokenKeywordIPv6Var
	TokenKeywordPortVar
	TokenKeywordVar
	TokenStringValue
	TokenIpv4Value
	TokenIpv4NetworkValue
	TokenIpv6Value
	TokenIpv6NetworkValue
	TokenVarUse
)

// Regular expressions for different patterns.
const (
	ipv4Pattern      = `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	ipv4RangePattern = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/(3[0-2]|[12]?[0-9])$`
	ipv6Pattern      = `^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4})|(([0-9a-fA-F]{1,4}:){1,6}:([0-9a-fA-F]{1,4}:){1,6})|(([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,6})|(([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,7})|(([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,8})|(([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,9})|(([0-9a-fA-F]{1,4}:){1}(:[0-9a-fA-F]{1,4}){1,10})|(:((:[0-9a-fA-F]{1,4}){1,11}))$`
)

// ConfigLexer represents a lexer for configuration files.
type ConfigLexer struct {
	input      string
	pos        int
	currLine   int
	curr       rune
	reader     *bufio.Reader
	file       *os.File
	CurrTokens []Token
}

// Custom error type for end-of-file.
type currCharError string

func (e *currCharError) Error() string {
	return fmt.Sprintf("currCharError: %s", string(*e))
}

// NextToken returns the next token from the input.
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

		if unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '$' || l.curr == '/' || l.curr == '_' || l.curr == '.' || l.curr == '[' || l.curr == ']' || l.curr == ',' {
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

// isSpecial checks if a character is special.
func (l *ConfigLexer) isSpecial(ch rune) bool {
	specials := []rune{'{', '}', '$'}
	for _, s := range specials {
		if ch == s {
			return true
		}
	}
	return false
}

// readChar reads the next character from the input.
func (l *ConfigLexer) readChar() {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		l.curr = eof
		return
	}
	l.curr = ch
	l.pos++
}

// skipWhitespace skips whitespace characters.
func (l *ConfigLexer) skipWhitespace() {
	for unicode.IsSpace(l.curr) {
		l.readChar()
	}
}

// skipComment skips comments.
func (l *ConfigLexer) skipComment() {
	for l.curr != '\n' && l.curr != eof {
		l.readChar()
	}
}

// parseIdentifier parses an identifier or keyword.
func (l *ConfigLexer) parseIdentifier() Token {
	startPos := l.pos - 1
	havePoint, haveBias, haveColon := false, false, false

	var haveDollar = l.curr == '$'
	if haveDollar {
		l.readChar()
		for unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '_' {
			l.readChar()
		}
		identifier := l.input[startPos+1 : l.pos-1]
		return Token{Type: TokenVarUse, Value: identifier}
	}

	for unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '.' || l.curr == '/' || l.curr == ':' || l.curr == '_' || l.curr == '[' || l.curr == ']' || l.curr == ',' {
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

// LexerInterface represents the interface for a lexer.
type LexerInterface interface {
	Tokenize() []Token
}

// NewConfigLexer creates a new lexer for the given file path.
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

// GetAllToken tokenizes the input and returns the sequence of tokens.
func (l *ConfigLexer) GetAllToken() []Token {
	for {
		myToken, _ := l.NextToken()
		if myToken.Type == TokenEOF {
			break
		}
	}
	return l.CurrTokens
}
