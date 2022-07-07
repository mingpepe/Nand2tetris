package tokenizer

const (
	KEYWORD int = iota
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

const (
	CLASS int = iota
	METHOD
	FUNCTION
	CONSTRUCTOR
	INT
	BOLLEAN
	CHAR
	VOID
	VAR
	STATIC
	FIELD
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN
	TRUE
	FALSE
	NULL
	THIS
)

type Tokenizer struct {
}

func (t *Tokenizer) hasMoreTokens() bool {
	return true
}

func (t *Tokenizer) advance() {

}

func (t *Tokenizer) tokenType() int {
	return KEYWORD
}

func (t *Tokenizer) keyword() int {
	return CLASS
}

func (t *Tokenizer) symbol() byte {
	return 0
}

func (t *Tokenizer) identifier() string {
	return ""
}

func (t *Tokenizer) intVal() int {
	return 0
}

func (t *Tokenizer) stringVal() string {
	return ""
}
