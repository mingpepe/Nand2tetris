// Project11
package compiler

import (
	"fmt"
	"io"
	"log"
	"runtime/debug"
)

type CompilationEngineVM struct {
	tokenizer      *Tokenizer
	vmWriter       *VMWriter
	symbolTable    *SymbolTable
	className      string
	labelPairCnt   int
	functionName   string
	subroutineType string
}

func NewCompilationEngineVM(reader io.Reader, writer io.Writer) *CompilationEngineVM {
	tokenizer := NewTokenizer(reader)
	tokenizer.Parse()
	tokenizer.Advance()
	vmWriter := NewVMWriter(writer)
	symbolTable := NewSymbolTable()

	c := &CompilationEngineVM{}
	c.tokenizer = tokenizer
	c.vmWriter = vmWriter
	c.symbolTable = symbolTable
	return c
}

func (e *CompilationEngineVM) mustHaveTokeType(tokenType string) {
	tok := e.tokenizer.TokenType()
	if tok != tokenType {
		debug.PrintStack()
		msg := fmt.Sprintf("unexpected token type(expected : %s, actual : %s)", tokenType, tok)
		log.Fatal(msg)
	}
}

func (e *CompilationEngineVM) mustHaveKeyword(keyword string) {
	e.mustHaveTokeType(KEYWORD)
	key := e.tokenizer.Keyword()
	if key != keyword {
		debug.PrintStack()
		msg := fmt.Sprintf("unexpected keyowrd(expected : %s, actaul : %s)", keyword, key)
		log.Fatal(msg)
	}
}

func (e *CompilationEngineVM) mustHaveSymbol(symbol byte) {
	e.mustHaveTokeType(SYMBOL)
	sym := e.tokenizer.Symbol()
	if sym != symbol {
		debug.PrintStack()
		msg := fmt.Sprintf("unexpected symbol(expected : %s, actaul %s)", string(symbol), string(sym))
		log.Fatal(msg)
	}
}

func (e *CompilationEngineVM) handleKeyword(keyword string) {
	e.mustHaveKeyword(keyword)
	e.tokenizer.Advance()
}

func (e *CompilationEngineVM) handleSymbol(symbol byte) {
	e.mustHaveSymbol(symbol)
	e.tokenizer.Advance()
}

func (e *CompilationEngineVM) CompileClass() {
	e.mustHaveKeyword(CLASS)
	e.tokenizer.Advance()

	e.mustHaveTokeType(IDENTIFIER)
	id := e.tokenizer.Identifier()
	e.className = id
	e.tokenizer.Advance()

	e.handleSymbol('{')

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
	e.tokenizer.Advance()
}

func (e *CompilationEngineVM) CompileClassVarDec() {
	e.mustHaveTokeType(KEYWORD)
	keyword := e.tokenizer.Keyword()
	if keyword != "static" && keyword != "field" {
		log.Fatalf("unexpected class var kind: %s", keyword)
	}
	e.tokenizer.Advance()

	var varType = ""
	tokType := e.tokenizer.TokenType()
	if tokType == KEYWORD {
		// int, char, boolean
		varType = e.tokenizer.Keyword()
		e.tokenizer.Advance()

	} else if tokType == IDENTIFIER {
		// class name
		varType = e.tokenizer.Identifier()
		e.tokenizer.Advance()

	} else {
		log.Fatalf("expected type keyword or identifier, got %s", tokType)
	}
	if e.tokenizer.TokenType() != IDENTIFIER {
		log.Fatal("expected variable name")
	}

	varName := e.tokenizer.Identifier()
	e.tokenizer.Advance()
	if keyword == "static" {
		e.symbolTable.Define(varName, varType, SYMBOL_STATIC)
	} else {
		e.symbolTable.Define(varName, varType, SYMBOL_FIELD)
	}

	if e.tokenizer.TokenType() == SYMBOL {
		symbol := e.tokenizer.Symbol()
		if symbol == ',' {
			e.handleSymbol(',')
			for {
				varName = e.tokenizer.Identifier()
				e.tokenizer.Advance()
				if keyword == "static" {
					e.symbolTable.Define(varName, varType, SYMBOL_STATIC)
				} else {
					e.symbolTable.Define(varName, varType, SYMBOL_FIELD)
				}
				if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == ';' {
					e.handleSymbol(';')
					break
				}
				e.handleSymbol(',')
			}
		} else if symbol == ';' {
			e.handleSymbol(';')
		} else {
			tmp := fmt.Sprintf("unexpected symbol : %s", string(symbol))
			log.Fatal(tmp)
		}
	} else {
		tokenType := e.tokenizer.TokenType()
		tmp := fmt.Sprintf("unexpected token type : %s", tokenType)
		log.Fatal(tmp)
	}
}

func (e *CompilationEngineVM) CompileSubroutineDec() {
	e.symbolTable.StartSubroutine()
	e.mustHaveTokeType(KEYWORD)
	key := e.tokenizer.Keyword()

	switch key {
	case CONSTRUCTOR:
		e.handleKeyword(CONSTRUCTOR)
		e.subroutineType = CONSTRUCTOR
	case FUNCTION:
		e.handleKeyword(FUNCTION)
		e.subroutineType = FUNCTION
	case METHOD:
		e.handleKeyword(METHOD)
		e.subroutineType = METHOD
	default:
		tmp := fmt.Sprintf("unexpected keyword : %s", key)
		log.Fatal(tmp)
	}

	// Return type, we do not need to generate vm code here
	tokenType := e.tokenizer.TokenType()
	if tokenType == KEYWORD {
		// void, int...
		key := e.tokenizer.CurrentToken()
		e.handleKeyword(key)
	} else if tokenType == IDENTIFIER {
		// User defined class
		//e.writeIdentifier()
		e.mustHaveTokeType(IDENTIFIER)
		e.tokenizer.Advance()
	} else {
		tmp := fmt.Sprintf("unexpected token type : %s", tokenType)
		log.Fatal(tmp)
	}

	// Subroutine name
	e.mustHaveTokeType(IDENTIFIER)
	id := e.tokenizer.Identifier()
	e.tokenizer.Advance()
	e.functionName = e.className + "." + id

	e.handleSymbol('(')
	e.CompileParameterList()
	e.handleSymbol(')')
	e.CompileSubroutineBody()
}

func (e *CompilationEngineVM) CompileParameterList() int {
	// Leave ')' to be handled in CompileSubroutineDec
	if e.tokenizer.TokenType() == SYMBOL {
		// No parameter
		e.mustHaveSymbol(')')
		return 0
	}

	cnt := 0
	for {
		var paramType, paramName string

		if e.tokenizer.TokenType() == KEYWORD {
			paramType = e.tokenizer.Keyword()
		} else if e.tokenizer.TokenType() == IDENTIFIER {
			paramType = e.tokenizer.Keyword()
		}
		e.tokenizer.Advance()

		e.mustHaveTokeType(IDENTIFIER)
		paramName = e.tokenizer.Identifier()
		e.tokenizer.Advance()
		cnt++

		e.symbolTable.Define(paramName, paramType, SYMBOL_ARG)
		if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == ',' {
			e.handleSymbol(',')
		} else {
			break
		}
	}
	return cnt
}

func (e *CompilationEngineVM) CompileSubroutineBody() {
	e.handleSymbol('{')
	for {
		if e.tokenizer.TokenType() == KEYWORD && e.tokenizer.Keyword() == VAR {
			e.CompileVarDec()
		} else {
			break
		}
	}

	nLocals := e.symbolTable.VarCount(SYMBOL_VAR)
	e.vmWriter.WriteFunction(e.functionName, nLocals)

	if e.subroutineType == CONSTRUCTOR {
		fieldCount := e.symbolTable.VarCount(SYMBOL_FIELD)
		e.vmWriter.WritePush("constant", fieldCount)
		e.vmWriter.WriteCall("Memory.alloc", 1)
		e.vmWriter.WritePop("pointer", 0)
	} else if e.subroutineType == METHOD {
		e.vmWriter.WritePush("argument", 0)
		e.vmWriter.WritePop("pointer", 0)
	}

	e.CompileStatements()
	e.handleSymbol('}')
}

func (e *CompilationEngineVM) CompileVarDec() {
	e.handleKeyword(VAR)

	var typeName string
	if e.tokenizer.TokenType() == IDENTIFIER {
		e.mustHaveTokeType(IDENTIFIER)
		typeName = e.tokenizer.Identifier()
		e.tokenizer.Advance()
	} else {
		// Fixme : not all keyword are allowed here
		keyword := e.tokenizer.Keyword()
		typeName = keyword
		e.handleKeyword(keyword)
	}

	for {
		e.mustHaveTokeType(IDENTIFIER)
		id := e.tokenizer.Identifier()
		e.symbolTable.Define(id, typeName, SYMBOL_VAR)
		e.tokenizer.Advance()

		e.mustHaveTokeType(SYMBOL)
		symbol := e.tokenizer.Symbol()
		if symbol == ',' {
			e.handleSymbol(symbol)
			continue
		} else if symbol == ';' {
			e.handleSymbol(symbol)
			break
		} else {
			debug.PrintStack()
			log.Fatalf("unexpected symbol : %s", string(symbol))
		}
	}
}

func (e *CompilationEngineVM) CompileStatements() {
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
}

func (e *CompilationEngineVM) CompileLet() {
	e.handleKeyword(LET)
	e.mustHaveTokeType(IDENTIFIER)
	varName := e.tokenizer.Identifier()
	e.vmWriter.WriteComment(fmt.Sprintf("CompileLet %s", varName))
	e.tokenizer.Advance()

	isArray := false
	if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == '[' {
		isArray = true
		e.handleSymbol('[')
		kind := e.symbolTable.KindOf(varName)
		segment := KindToSegment(kind)
		index := e.symbolTable.IndexOf(varName)
		e.vmWriter.WritePush(segment, index)
		e.CompileExpression()
		e.handleSymbol(']')

		e.vmWriter.WriteArithmetic("add")
	}
	e.handleSymbol('=')
	e.CompileExpression()
	if isArray {
		// store RHS
		e.vmWriter.WritePop("temp", 0)
		// set THAT to address
		e.vmWriter.WritePop("pointer", 1)
		// restore RHS
		e.vmWriter.WritePush("temp", 0)
		e.vmWriter.WritePop("that", 0)
	} else {
		kind := e.symbolTable.KindOf(varName)
		segment := KindToSegment(kind)
		idx := e.symbolTable.IndexOf(varName)
		e.vmWriter.WritePop(segment, idx)
	}
	e.handleSymbol(';')
}

func (e *CompilationEngineVM) CompileIf() {
	idx := e.labelPairCnt
	e.labelPairCnt++
	L1 := fmt.Sprintf("IF_L1%d", idx)
	L2 := fmt.Sprintf("IF_L2%d", idx)
	L3 := fmt.Sprintf("IF_L3%d", idx)

	e.handleKeyword(IF)
	e.handleSymbol('(')
	e.CompileExpression()
	e.handleSymbol(')')

	// Not instruction in VM is bitwise not. e.g. not 1 => -2
	// Use 3 labels to handle it
	e.vmWriter.WriteIf(L1)
	e.vmWriter.WriteGoTo(L2)
	e.vmWriter.WriteLabel(L1)
	e.handleSymbol('{')
	e.CompileStatements()
	e.handleSymbol('}')
	e.vmWriter.WriteGoTo(L3)

	e.vmWriter.WriteLabel(L2)
	if e.tokenizer.TokenType() == KEYWORD && e.tokenizer.Keyword() == ELSE {
		e.handleKeyword(ELSE)
		e.handleSymbol('{')
		e.CompileStatements()
		e.handleSymbol('}')
	}
	e.vmWriter.WriteLabel(L3)

}

func (e *CompilationEngineVM) CompileWhile() {
	startLabel := fmt.Sprintf("WHILE_START%d", e.labelPairCnt)
	endLabel := fmt.Sprintf("WHILE_END%d", e.labelPairCnt)
	e.labelPairCnt++
	e.handleKeyword(WHILE)

	e.vmWriter.WriteLabel(startLabel)
	e.handleSymbol('(')
	e.CompileExpression()
	e.handleSymbol(')')

	e.vmWriter.WriteArithmetic("not")
	e.vmWriter.WriteIf(endLabel)

	e.handleSymbol('{')
	e.CompileStatements()
	e.handleSymbol('}')

	e.vmWriter.WriteGoTo(startLabel)
	e.vmWriter.WriteLabel(endLabel)
}

func (e *CompilationEngineVM) CompileDo() {
	e.handleKeyword(DO)
	e.mustHaveTokeType(IDENTIFIER)
	firstName := e.tokenizer.Identifier()
	e.tokenizer.Advance()
	var fullName string
	nArgs := 0

	e.mustHaveTokeType(SYMBOL)

	if e.tokenizer.Symbol() == '.' {
		e.tokenizer.Advance()

		e.mustHaveTokeType(IDENTIFIER)
		secondName := e.tokenizer.Identifier()
		e.tokenizer.Advance()

		kind := e.symbolTable.KindOf(firstName)
		// method call
		if kind != SYMBOL_NONE {
			// push object reference
			segment := KindToSegment(kind)
			e.vmWriter.WritePush(segment, e.symbolTable.IndexOf(firstName))
			// fullName = ClassName.method
			fullName = e.symbolTable.TypeOf(firstName) + "." + secondName
			nArgs = 1
		} else {
			fullName = firstName + "." + secondName
		}
	} else {
		fullName = e.className + "." + firstName
		// push this
		e.vmWriter.WritePush("pointer", 0)
		nArgs = 1
	}

	e.handleSymbol('(')
	nArgs += e.CompileExpressionList()
	e.handleSymbol(')')
	e.vmWriter.WriteCall(fullName, nArgs)
	e.vmWriter.WritePop("temp", 0)
	e.handleSymbol(';')
}

func (e *CompilationEngineVM) CompileReturn() {
	e.handleKeyword(RETURN)
	if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == ';' {
		// Empty
		e.vmWriter.WritePush("constant", 0)
	} else {
		e.CompileExpression()
	}
	e.handleSymbol(';')
	e.vmWriter.WriteReturn()
}

func (e *CompilationEngineVM) CompileExpression() {
	e.CompileTerm()
	for {
		if e.tokenizer.TokenType() == SYMBOL {
			symbol := e.tokenizer.Symbol()
			match := false
			ops := []byte{'+', '-', '*', '/', '&', '|', '<', '>', '='}
			for _, op := range ops {
				if symbol == op {
					match = true
					e.tokenizer.Advance()
					e.CompileTerm() // right side
					switch symbol {
					case '+':
						e.vmWriter.WriteArithmetic("add")
					case '-':
						e.vmWriter.WriteArithmetic("sub")
					case '*':
						e.vmWriter.WriteCall("Math.multiply", 2)
					case '/':
						e.vmWriter.WriteCall("Math.divide", 2)
					case '&':
						e.vmWriter.WriteArithmetic("and")
					case '|':
						e.vmWriter.WriteArithmetic("or")
					case '<':
						e.vmWriter.WriteArithmetic("lt")
					case '>':
						e.vmWriter.WriteArithmetic("gt")
					case '=':
						e.vmWriter.WriteArithmetic("eq")
					}
					break
				}
			}
			if match {
				continue
			}
		}
		break
	}
}

func (e *CompilationEngineVM) CompileTerm() {
	if e.tokenizer.TokenType() == STRING_CONST {
		str := e.tokenizer.StringVal()
		e.vmWriter.WritePush("constant", len(str))
		e.vmWriter.WriteCall("String.new", 1)
		for _, ch := range str {
			e.vmWriter.WritePush("constant", int(ch))
			e.vmWriter.WriteCall("String.appendChar", 2)
		}
		e.tokenizer.Advance()
	} else if e.tokenizer.TokenType() == INT_CONST {
		val := e.tokenizer.IntVal()
		e.vmWriter.WritePush("constant", val)
		e.tokenizer.Advance()
	} else if e.tokenizer.TokenType() == IDENTIFIER {
		e.mustHaveTokeType(IDENTIFIER)
		name := e.tokenizer.Identifier()
		e.tokenizer.Advance()

		if e.tokenizer.TokenType() == SYMBOL {
			if e.tokenizer.Symbol() == '[' {
				kind := e.symbolTable.KindOf(name)
				segment := KindToSegment(kind)
				index := e.symbolTable.IndexOf(name)
				e.vmWriter.WritePush(segment, index)
				e.handleSymbol('[')
				e.CompileExpression()
				e.handleSymbol(']')

				e.vmWriter.WriteArithmetic("add")
				e.vmWriter.WritePop("pointer", 1)
				e.vmWriter.WritePush("that", 0)

			} else if e.tokenizer.Symbol() == '.' || e.tokenizer.Symbol() == '(' {
				fullName := ""
				nArgs := 0
				if e.tokenizer.Symbol() == '(' {
					fullName = e.className + "." + name
				} else {
					e.handleSymbol('.')
					subName := e.tokenizer.Identifier()
					e.tokenizer.Advance()

					if kind := e.symbolTable.KindOf(name); kind != SYMBOL_NONE {
						// varName.methodName
						segment := KindToSegment(kind)
						index := e.symbolTable.IndexOf(name)
						e.vmWriter.WritePush(segment, index)
						fullName = e.symbolTable.TypeOf(name) + "." + subName
						nArgs++ // add 'this'
					} else {
						// className.subroutineName
						fullName = name + "." + subName
					}
				}

				e.handleSymbol('(')
				nArgs += e.CompileExpressionList()
				e.handleSymbol(')')
				e.vmWriter.WriteCall(fullName, nArgs)
			} else {
				// Simple var
				kind := e.symbolTable.KindOf(name)
				segment := KindToSegment(kind)
				index := e.symbolTable.IndexOf(name)
				if segment == "argument" && e.subroutineType == METHOD {
					index++
				}
				e.vmWriter.WritePush(segment, index)
			}
		}
	} else if e.tokenizer.TokenType() == KEYWORD {
		switch e.tokenizer.Keyword() {
		case TRUE:
			e.vmWriter.WritePush("constant", 0)
			e.vmWriter.WriteArithmetic("not")
		case FALSE, NULL:
			e.vmWriter.WritePush("constant", 0)
		case THIS:
			e.vmWriter.WritePush("pointer", 0)
		}
		e.tokenizer.Advance()
	} else if e.tokenizer.TokenType() == SYMBOL {
		symbol := e.tokenizer.Symbol()
		if symbol == '(' {
			e.handleSymbol('(')
			e.CompileExpression()
			e.handleSymbol(')')
		} else if symbol == '-' || symbol == '~' {
			e.handleSymbol(symbol)
			e.CompileTerm()
			if symbol == '-' {
				e.vmWriter.WriteArithmetic("neg")
			} else {
				e.vmWriter.WriteArithmetic("not")
			}
		} else {
			tmp := fmt.Sprintf("unexpected symbol : %s", string(symbol))
			log.Fatal(tmp)
		}

	} else {
		tokenType := e.tokenizer.TokenType()
		tmp := fmt.Sprintf("unexpected token type : %s", tokenType)
		log.Fatal(tmp)
	}
}

func (e *CompilationEngineVM) CompileExpressionList() (nArgs int) {
	nArgs = 0
	if e.tokenizer.TokenType() == SYMBOL && e.tokenizer.Symbol() == ')' {
		goto END
	}
	for {
		nArgs++
		e.CompileExpression()
		if e.tokenizer.TokenType() == SYMBOL {
			if e.tokenizer.Symbol() == ',' {
				e.handleSymbol(',')
				continue
			}
		}
		break
	}
END:
	return nArgs
}
