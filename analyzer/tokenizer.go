package analyzer

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// Token type
const (
	KEYWORD      = "keyword"
	SYMBOL       = "symbol"
	IDENTIFIER   = "identifier"
	INT_CONST    = "integerConstant"
	STRING_CONST = "stringConstant"
)

// Keyword
const (
	CLASS       = "class"
	METHOD      = "method"
	FUNCTION    = "function"
	CONSTRUCTOR = "constructor"
	INT         = "int"
	BOOLEAN     = "boolean"
	CHAR        = "char"
	VOID        = "void"
	VAR         = "var"
	STATIC      = "static"
	FIELD       = "field"
	LET         = "let"
	DO          = "do"
	IF          = "if"
	ELSE        = "else"
	WHILE       = "while"
	RETURN      = "return"
	TRUE        = "true"
	FALSE       = "false"
	NULL        = "null"
	THIS        = "this"
)

var keywords = []string{CLASS, METHOD, FUNCTION,
	CONSTRUCTOR, INT, BOOLEAN,
	CHAR, VOID, VAR, STATIC,
	FIELD, LET, DO, IF,
	ELSE, WHILE, RETURN, TRUE,
	FALSE, NULL, THIS}

const symbols = "{}()[].,;+-*/&|<>=~"

const sep = " \t\r\n"

type Tokenizer struct {
	reader io.Reader
	tokens []string
	ptr    int
}

func NewTokenizer(reader io.Reader) *Tokenizer {
	t := &Tokenizer{}
	t.reader = reader
	t.tokens = make([]string, 0)
	t.ptr = -1
	return t
}

func (t *Tokenizer) Parse() {
	buf := make([]rune, 0)
	ptr := 0
	scanner := bufio.NewScanner(t.reader)
	multi_line_comments := false
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "/*") {
			if strings.HasSuffix(line, "*/") {
				continue
			} else {
				multi_line_comments = true
			}
		}

		if multi_line_comments {
			// Assume that it's the end of line
			idx := strings.Index(line, "*/")
			if idx != -1 {
				multi_line_comments = false
			}
			continue
		}

		idx := strings.Index(line, "//")
		if idx != -1 {
			line = line[:idx]
		}

		// Does not handle cross line string
		string_start_flag := false
		for _, b := range line {
			if string_start_flag {
				buf = append(buf, b)
				ptr++
				if b == '"' {
					string_start_flag = false
					s := string(buf[:ptr])
					t.tokens = append(t.tokens, s)
					ptr = 0
					buf = make([]rune, 0)
				}
			} else {
				if strings.Contains(sep, string(b)) {
					if ptr != 0 {
						s := string(buf[:ptr])
						t.tokens = append(t.tokens, s)
						ptr = 0
						buf = make([]rune, 0)
					}
				} else if strings.Contains(symbols, string(b)) {
					if ptr != 0 {
						s := string(buf[:ptr])
						t.tokens = append(t.tokens, s)
						ptr = 0
						buf = make([]rune, 0)
					}
					t.tokens = append(t.tokens, string(b))
				} else {
					buf = append(buf, b)
					ptr++
					if b == '"' {
						string_start_flag = true
					}
				}
			}

		}
		if ptr != 0 {
			s := string(buf[:ptr])
			t.tokens = append(t.tokens, s)
			ptr = 0
			buf = make([]rune, 0)
		}
	}
}

func (t *Tokenizer) HasMoreTokens() bool {
	return t.ptr < len(t.tokens)-1
}

func (t *Tokenizer) Advance() {
	t.ptr++
}

func (t *Tokenizer) CurrentToken() string {
	token := t.tokens[t.ptr]
	if t.TokenType() == STRING_CONST {
		length := len(token)
		token = token[1 : length-1]
	} else if token == "<" {
		token = "&lt;"
	} else if token == ">" {
		token = "&gt;"
	} else if token == "&" {
		token = "&amp;"
	}
	return token
}

func (t *Tokenizer) TokenType() string {
	token := t.tokens[t.ptr]
	if strings.HasPrefix(token, "\"") {
		return STRING_CONST
	}

	if strings.Contains(symbols, token) {
		return SYMBOL
	}

	for _, keyword := range keywords {
		if token == keyword {
			return KEYWORD
		}
	}

	_, err := strconv.Atoi(token)
	if err == nil {
		return INT_CONST
	}
	return IDENTIFIER
}

func (t *Tokenizer) Keyword() string {
	return t.CurrentToken()
}

func (t *Tokenizer) Symbol() byte {
	token := t.CurrentToken()
	switch token {
	case "&lt;":
		return '<'
	case "&gt;":
		return '>'
	case "&amp;":
		return '&'
	}
	return token[0]
}

func (t *Tokenizer) Identifier() string {
	return t.CurrentToken()
}

func (t *Tokenizer) IntVal() int {
	s, _ := strconv.Atoi(t.CurrentToken())
	return s
}

func (t *Tokenizer) StringVal() string {
	return t.CurrentToken()
}
