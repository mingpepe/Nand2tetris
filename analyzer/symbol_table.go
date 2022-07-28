package analyzer

type SymbolTable struct {
}

const (
	SYMBOL_STATIC int = iota
	SYMBOL_FIELD
	SYMBOL_ARG
	SYMBOL_VAR

	SYMBOL_NONE
)

func NewSymbolTable() *SymbolTable {
	s := SymbolTable{}
	return &s
}

func (s *SymbolTable) StartSubroutine() {

}

func (s *SymbolTable) Define(name, _type string, kind int) {

}

func (s *SymbolTable) VarCount(kind int) int {
	return 0
}

func (s *SymbolTable) KingOf(name string) int {
	return SYMBOL_NONE
}

func (s *SymbolTable) TypeOf(name string) string {
	return ""
}

func (s *SymbolTable) IndexOf(name string) int {
	return 0
}
