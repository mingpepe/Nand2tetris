package compiler

import (
	"log"
)

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

func KindToSegment(kind KIND) string {
	switch kind {
	case SYMBOL_STATIC:
		return "static"
	case SYMBOL_FIELD:
		return "this"
	case SYMBOL_ARG:
		return "argument"
	case SYMBOL_VAR:
		return "local"
	default:
		log.Fatalf("unknown kind: %d", kind)
		return ""
	}
}

func NewSymbolTable() *SymbolTable {
	s := SymbolTable{}
	s.table = make(map[string]Variable)
	s.subRoutinetable = make(map[string]Variable)
	s.kindCountTable = make(map[KIND]int)
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
	case SYMBOL_STATIC, SYMBOL_FIELD:
		s.table[name] = variable
	case SYMBOL_ARG, SYMBOL_VAR:
		s.subRoutinetable[name] = variable
	}
	s.kindCountTable[kind]++
}

func (s *SymbolTable) VarCount(kind KIND) int {
	return s.kindCountTable[kind]
}

func (s *SymbolTable) KindOf(name string) KIND {
	return s.lookUp(name).kind
}

func (s *SymbolTable) TypeOf(name string) string {
	return s.lookUp(name).varType
}

func (s *SymbolTable) IndexOf(name string) int {
	return s.lookUp(name).index
}

func (s *SymbolTable) lookUp(name string) Variable {
	if v, exist := s.table[name]; exist {
		return v
	}
	if v, exist := s.subRoutinetable[name]; exist {
		return v
	}
	return Variable{kind: SYMBOL_NONE}
}
