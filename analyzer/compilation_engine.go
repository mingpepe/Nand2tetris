package analyzer

import "io"

type CompilationEngine struct {
	tokenizer Tokenizer
	writer    io.Writer
}

func NewCompilationEngine(tokenizer Tokenizer, writer io.Writer) *CompilationEngine {
	c := &CompilationEngine{}
	c.tokenizer = tokenizer
	return c
}

func (e *CompilationEngine) compileClass() {

}

func (e *CompilationEngine) compileClassVarDec() {

}

func (e *CompilationEngine) compileSubroutineDec() {

}

func (e *CompilationEngine) compileParameterList() {

}

func (e *CompilationEngine) compileSubroutineBody() {

}

func (e *CompilationEngine) compileVarDec() {

}

func (e *CompilationEngine) compileStatements() {

}

func (e *CompilationEngine) compileLet() {

}

func (e *CompilationEngine) compileIf() {

}

func (e *CompilationEngine) compileWhile() {

}

func (e *CompilationEngine) compileDo() {

}

func (e *CompilationEngine) compileReturn() {

}
