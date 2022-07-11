package analyzer

import (
	"fmt"
	"io"
	"log"
	"runtime/debug"
)

type CompilationEngine struct {
	tokenizer *Tokenizer
	writer    io.Writer
}

func NewCompilationEngine(tokenizer *Tokenizer, writer io.Writer) *CompilationEngine {
	c := &CompilationEngine{}
	c.tokenizer = tokenizer
	c.writer = writer
	return c
}

func (e *CompilationEngine) mustHaveTokeType(tokenType string) {
	tok := e.tokenizer.TokenType()
	if tok != tokenType {
		debug.PrintStack()
		msg := fmt.Sprintf("unexpected token type(expected : %s, actual : %s)", tokenType, tok)
		log.Fatal(msg)
	}
}

func (e *CompilationEngine) mustHaveKeyword(keyword string) {
	e.mustHaveTokeType(KEYWORD)
	key := e.tokenizer.Keyword()
	if key != keyword {
		debug.PrintStack()
		msg := fmt.Sprintf("unexpected keyowrd(expected : %s, actaul : %s)", keyword, key)
		log.Fatal(msg)
	}
}

func (e *CompilationEngine) mustHaveSymol(symbol byte) {
	e.mustHaveTokeType(SYMBOL)
	sym := e.tokenizer.Symbol()
	if sym != symbol {
		debug.PrintStack()
		msg := fmt.Sprintf("unexpected symbol(expected : %s, actaul %s)", string(symbol), string(sym))
		log.Fatal(msg)
	}
}

func (e *CompilationEngine) writeOutput(content string) {
	e.writer.Write([]byte(content))
}

func (e *CompilationEngine) writeOutputLine(content string) {
	e.writeOutput(content + "\n")
}

func (e *CompilationEngine) writeKeyword(keyword string) {
	e.mustHaveKeyword(keyword)
	tmp := fmt.Sprintf("<keyword>%s</keyword>", keyword)
	e.writeOutputLine(tmp)
	e.tokenizer.Advance()
}

func (e *CompilationEngine) writeIdentifier() {
	e.mustHaveTokeType(IDENTIFIER)
	id := e.tokenizer.Identifier()
	tmp := fmt.Sprintf("<identifier>%s</identifier>", id)
	e.writeOutputLine(tmp)
	e.tokenizer.Advance()
}

func (e *CompilationEngine) writeSymbol(symbol byte) {
	e.mustHaveSymol(symbol)
	tmp := fmt.Sprintf("<symbol>%s</symbol>", string(symbol))
	e.writeOutputLine(tmp)
	e.tokenizer.Advance()
}

func (e *CompilationEngine) CompileClass() {
	e.writeOutputLine("<class>")
	e.writeKeyword(CLASS)
	e.writeIdentifier()
	e.writeSymbol('{')

	for {
		if e.tokenizer.TokenType() == KEYWORD {
			keyword := e.tokenizer.Keyword()
			if keyword == STATIC || keyword == FIELD {
				e.CompileClassVarDec()
				continue
			}
		}
		break
	}
	for {
		if e.tokenizer.TokenType() == KEYWORD {
			keyword := e.tokenizer.Keyword()
			if keyword == FUNCTION || keyword == CONSTRUCTOR || keyword == METHOD {
				e.CompileSubroutineDec()
				continue
			}
		}
		break
	}
	e.writeSymbol('}')
	e.writeOutputLine("</class>")
}

func (e *CompilationEngine) CompileClassVarDec() {
	e.writeOutputLine("<classVarDec>")
	keyword := e.tokenizer.Keyword()
	e.writeKeyword(keyword)
	tokenType := e.tokenizer.TokenType()
	if tokenType == KEYWORD {
		// FIXME : not all keywords are allowed
		keyword = e.tokenizer.Keyword()
		e.writeKeyword(keyword)
	} else if tokenType == IDENTIFIER {
		e.writeIdentifier()
	} else {
		tmp := fmt.Sprintf("unexpected token type : %s", tokenType)
		log.Fatal(tmp)
	}

	e.writeIdentifier()
	if e.tokenizer.TokenType() == SYMBOL {
		symbol := e.tokenizer.Symbol()
		if symbol == ',' {
			e.writeSymbol(',')
			for {
				e.writeIdentifier()
				if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == ';' {
					e.writeSymbol(';')
					break
				}
				e.writeSymbol(',')
			}
		} else if symbol == ';' {
			e.writeSymbol(';')
		} else {
			tmp := fmt.Sprintf("unexpected symbol : %s", string(symbol))
			log.Fatal(tmp)
		}
	} else {
		tokenType := e.tokenizer.TokenType()
		tmp := fmt.Sprintf("unexpected token type : %s", tokenType)
		log.Fatal(tmp)
	}
	e.writeOutputLine("</classVarDec>")
}

func (e *CompilationEngine) CompileSubroutineDec() {
	e.writeOutputLine("<subroutineDec>")
	e.mustHaveTokeType(KEYWORD)
	key := e.tokenizer.Keyword()

	switch key {
	case CONSTRUCTOR:
		e.writeKeyword(CONSTRUCTOR)
	case FUNCTION:
		e.writeKeyword(FUNCTION)
	case METHOD:
		e.writeKeyword(METHOD)
	default:
		tmp := fmt.Sprintf("unexpected keyword : %s", key)
		log.Fatal(tmp)
	}

	tokenType := e.tokenizer.TokenType()
	if tokenType == KEYWORD {
		key := e.tokenizer.CurrentToken()
		e.writeKeyword(key)
	} else if tokenType == IDENTIFIER {
		e.writeIdentifier()
	} else {
		tmp := fmt.Sprintf("unexpected token type : %s", tokenType)
		log.Fatal(tmp)
	}

	e.writeIdentifier()
	e.writeSymbol('(')
	e.CompileParameterList()
	e.writeSymbol(')')
	e.CompileSubroutineBody()
	e.writeOutputLine("</subroutineDec>")
}

func (e *CompilationEngine) CompileParameterList() {
	e.writeOutputLine("<parameterList>")
	if e.tokenizer.TokenType() == SYMBOL {
		// No parameter
		e.mustHaveSymol(')')
	} else {
		for {
			if e.tokenizer.TokenType() == KEYWORD {
				// FIXME : not all keywords are allowed
				keyowrd := e.tokenizer.Keyword()
				e.writeKeyword(keyowrd)
			} else if e.tokenizer.TokenType() == IDENTIFIER {
				e.writeIdentifier()
			}
			e.writeIdentifier()
			if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == ',' {
				e.writeSymbol(',')
			} else {
				break
			}
		}

	}
	e.writeOutputLine("</parameterList>")
}

func (e *CompilationEngine) CompileSubroutineBody() {
	e.writeOutputLine("<subroutineBody>")
	e.writeSymbol('{')
	for {
		if e.tokenizer.TokenType() == KEYWORD && e.tokenizer.Keyword() == VAR {
			e.CompileVarDec()
		} else {
			break
		}
	}
	e.CompileStatements()
	e.writeSymbol('}')
	e.writeOutputLine("</subroutineBody>")
}

func (e *CompilationEngine) CompileVarDec() {
	e.writeOutputLine("<varDec>")
	e.writeKeyword(VAR)

	if e.tokenizer.TokenType() == IDENTIFIER {
		e.writeIdentifier()
	} else {
		// Fixme : not all keyword are allowed here
		keyword := e.tokenizer.Keyword()
		e.writeKeyword(keyword)
	}

	for {
		e.writeIdentifier()

		e.mustHaveTokeType(SYMBOL)
		symbol := e.tokenizer.Symbol()
		if symbol == ',' {
			e.writeSymbol(symbol)
			continue
		} else if symbol == ';' {
			e.writeSymbol(symbol)
			break
		} else {
			debug.PrintStack()
			log.Fatal(fmt.Sprintf("unexpected symbol : %s", string(symbol)))
		}
	}
	e.writeOutputLine("</varDec>")
}

func (e *CompilationEngine) CompileStatements() {
	e.writeOutputLine("<statements>")
	if e.tokenizer.TokenType() == KEYWORD {
		keep_going := true
		for keep_going {
			switch e.tokenizer.Keyword() {
			case LET:
				e.CompileLet()
			case IF:
				e.CompileIf()
			case WHILE:
				e.CompileWhile()
			case DO:
				e.CompileDo()
			case RETURN:
				e.CompileReturn()
				keep_going = false
			default:
				keep_going = false
			}
		}
	}
	e.writeOutputLine("</statements>")
}

func (e *CompilationEngine) CompileLet() {
	e.writeOutputLine("<letStatement>")
	e.writeKeyword(LET)
	e.writeIdentifier()
	if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == '[' {
		e.writeSymbol('[')
		e.CompileExpression()
		e.writeSymbol(']')
	}
	e.writeSymbol('=')
	e.CompileExpression()
	e.writeSymbol(';')
	e.writeOutputLine("</letStatement>")
}

func (e *CompilationEngine) CompileIf() {
	e.writeOutputLine("<ifStatement>")
	e.writeKeyword(IF)
	e.writeSymbol('(')
	e.CompileExpression()
	e.writeSymbol(')')
	e.writeSymbol('{')
	e.CompileStatements()
	e.writeSymbol('}')
	if e.tokenizer.TokenType() == KEYWORD && e.tokenizer.Keyword() == ELSE {
		e.writeKeyword(ELSE)
		e.writeSymbol('{')
		e.CompileStatements()
		e.writeSymbol('}')
	}
	e.writeOutputLine("</ifStatement>")
}

func (e *CompilationEngine) CompileWhile() {
	e.writeOutputLine("<whileStatement>")
	e.writeKeyword(WHILE)

	e.writeSymbol('(')
	e.CompileExpression()
	e.writeSymbol(')')
	e.writeSymbol('{')
	e.CompileStatements()
	e.writeSymbol('}')
	e.writeOutputLine("</whileStatement>")
}

func (e *CompilationEngine) CompileDo() {
	e.writeOutputLine("<doStatement>")
	e.writeKeyword(DO)

	for {
		e.writeIdentifier()
		e.mustHaveTokeType(SYMBOL)
		if e.tokenizer.Symbol() == '.' {
			e.writeSymbol('.')
		} else {
			break
		}
	}

	e.writeSymbol('(')
	e.CompileExpressionList()
	e.writeSymbol(')')
	e.writeSymbol(';')
	e.writeOutputLine("</doStatement>")
}

func (e *CompilationEngine) CompileReturn() {
	e.writeOutputLine("<returnStatement>")
	e.writeKeyword(RETURN)
	if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == ';' {
		// Empty
	} else {
		e.CompileExpression()
	}
	e.writeSymbol(';')
	e.writeOutputLine("</returnStatement>")
}

func (e *CompilationEngine) CompileExpression() {
	e.writeOutputLine("<expression>")
	for {
		e.CompileTerm()
		if e.tokenizer.TokenType() == SYMBOL {
			symbol := e.tokenizer.Symbol()
			match := false
			ops := []byte{'+', '-', '*', '/', '&', '|', '<', '>', '='}
			for _, op := range ops {
				if symbol == op {
					match = true
					tmp := op2xmlString(op)
					e.writeOutputLine(fmt.Sprintf("<symbol>%s</symbol>", tmp))
					e.tokenizer.Advance()
					break
				}
			}
			if match {
				continue
			}
		}
		break
	}
	e.writeOutputLine("</expression>")
}

func (e *CompilationEngine) CompileTerm() {
	e.writeOutputLine("<term>")
	if e.tokenizer.TokenType() == STRING_CONST {
		tmp := e.tokenizer.StringVal()
		e.writeOutputLine(fmt.Sprintf("<stringConstant>%s</stringConstant>", tmp))
		e.tokenizer.Advance()
	} else if e.tokenizer.TokenType() == INT_CONST {
		tmp := e.tokenizer.IntVal()
		e.writeOutputLine(fmt.Sprintf("<integerConstant>%d</integerConstant>", tmp))
		e.tokenizer.Advance()
	} else if e.tokenizer.TokenType() == IDENTIFIER {
		e.writeIdentifier()
		if e.tokenizer.TokenType() == SYMBOL {
			if e.tokenizer.Symbol() == '[' {
				e.writeSymbol('[')
				e.CompileExpression()
				e.writeSymbol(']')
			} else if e.tokenizer.Symbol() == '.' {
				e.writeSymbol('.')
				for {
					e.writeIdentifier()
					if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == '.' {
						e.writeSymbol('.')
						continue
					}
					break
				}
				e.writeSymbol('(')
				e.CompileExpressionList()
				e.writeSymbol(')')
			}
		}
	} else if e.tokenizer.TokenType() == KEYWORD {
		// FIXME : not all keywords are allowed
		keyword := e.tokenizer.Keyword()
		e.writeKeyword(keyword)
	} else if e.tokenizer.TokenType() == SYMBOL {
		symbol := e.tokenizer.Symbol()
		if symbol == '(' {
			e.writeSymbol('(')
			e.CompileExpression()
			e.writeSymbol(')')
		} else if symbol == '-' || symbol == '~' {
			e.writeSymbol(symbol)
			e.CompileTerm()
		} else {
			tmp := fmt.Sprintf("unexpected symbol : %s", string(symbol))
			log.Fatal(tmp)
		}

	} else {
		tokenType := e.tokenizer.TokenType()
		tmp := fmt.Sprintf("unexpected token type : %s", tokenType)
		log.Fatal(tmp)
	}
	e.writeOutputLine("</term>")
}

func (e *CompilationEngine) CompileExpressionList() {
	e.writeOutputLine("<expressionList>")
	if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == ')' {
		goto END
	}
	for {
		e.CompileExpression()
		if e.tokenizer.TokenType() == SYMBOL {
			if e.tokenizer.Symbol() == ',' {
				e.writeSymbol(',')
				continue
			}
		}
		break
	}
END:
	e.writeOutputLine("</expressionList>")
}

func op2xmlString(op byte) string {
	switch op {
	case '<':
		return "&lt;"
	case '>':
		return "&gt;"
	case '&':
		return "&amp;"
	}
	return string(op)
}
