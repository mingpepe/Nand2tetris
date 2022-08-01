package compiler

type KIND int
type SymbolTable struct {
	table          map[string]Variable
	kindCountTable map[KIND]int
}

type Variable struct {
	varType string
	kind    KIND
	index   int
}

const (
	SYMBOL_STATIC KIND = iota
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

func (s *SymbolTable) Define(name, _type string, kind KIND) {
	var variable Variable
	variable.varType = _type
	variable.kind = kind
	variable.index = s.kindCountTable[kind]
	s.table[name] = variable
	s.kindCountTable[kind]++
}

func (s *SymbolTable) VarCount(kind KIND) int {
	return s.kindCountTable[kind]
}

func (s *SymbolTable) KingOf(name string) KIND {
	return s.table[name].kind
}

func (s *SymbolTable) TypeOf(name string) string {
	return s.table[name].varType
}

func (s *SymbolTable) IndexOf(name string) int {
	return s.table[name].index
}
