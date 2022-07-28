package analyzer

type SymbolTable struct {
	table          map[string]Variable
	kindCountTable map[int]int
}

type Variable struct {
	varType string
	kind    int
	index   int
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
	s.kindCountTable[SYMBOL_STATIC] = 0
	s.kindCountTable[SYMBOL_FIELD] = 0
	s.kindCountTable[SYMBOL_ARG] = 0
	s.kindCountTable[SYMBOL_VAR] = 0
	return &s
}

func (s *SymbolTable) StartSubroutine() {

}

func (s *SymbolTable) Define(name, _type string, kind int) {
	var variable Variable
	variable.varType = _type
	variable.kind = kind
	variable.index = s.kindCountTable[kind]
	s.table[name] = variable
	s.kindCountTable[kind]++
}

func (s *SymbolTable) VarCount(kind int) int {
	return s.kindCountTable[kind]
}

func (s *SymbolTable) KingOf(name string) int {
	return s.table[name].kind
}

func (s *SymbolTable) TypeOf(name string) string {
	return s.table[name].varType
}

func (s *SymbolTable) IndexOf(name string) int {
	return s.table[name].index
}
