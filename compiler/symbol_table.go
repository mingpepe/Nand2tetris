package compiler

type KIND int
type SymbolTable struct {
	table           map[string]Variable
	subRoutinetable map[string]Variable
	kindCountTable  map[KIND]int
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
	s.table = make(map[string]Variable)
	s.subRoutinetable = make(map[string]Variable)
	s.kindCountTable[SYMBOL_STATIC] = 0
	s.kindCountTable[SYMBOL_FIELD] = 0
	s.kindCountTable[SYMBOL_ARG] = 0
	s.kindCountTable[SYMBOL_VAR] = 0
	return &s
}

func (s *SymbolTable) StartSubroutine() {
	s.kindCountTable[SYMBOL_ARG] = 0
	s.kindCountTable[SYMBOL_VAR] = 0
	s.subRoutinetable = make(map[string]Variable)
}

func (s *SymbolTable) Define(name, _type string, kind KIND) {
	var variable Variable
	variable.varType = _type
	variable.kind = kind
	variable.index = s.kindCountTable[kind]
	switch kind {
	case SYMBOL_STATIC:
	case SYMBOL_FIELD:
		s.table[name] = variable
	case SYMBOL_ARG:
	case SYMBOL_VAR:
		s.subRoutinetable[name] = variable
	}
	s.kindCountTable[kind]++
}

func (s *SymbolTable) VarCount(kind KIND) int {
	return s.kindCountTable[kind]
}

func (s *SymbolTable) KingOf(name string) KIND {
	return s.lookUp(name).kind
}

func (s *SymbolTable) TypeOf(name string) string {
	return s.lookUp(name).varType
}

func (s *SymbolTable) IndexOf(name string) int {
	return s.lookUp(name).index
}

func (s *SymbolTable) lookUp(name string) Variable {
	v, exist := s.table[name]
	if exist {
		return v
	}
	return s.subRoutinetable[name]
}
