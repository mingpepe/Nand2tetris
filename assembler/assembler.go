package assembler

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Assembler struct {
	destTable map[string]uint16
	compTable map[string]uint16
	jumpTable map[string]uint16
}

func New() *Assembler {
	a := &Assembler{}
	a.destTable = map[string]uint16{
		"M":   0b001000,
		"D":   0b010000,
		"DM":  0b011000,
		"A":   0b100000,
		"AM":  0b101000,
		"AD":  0b110000,
		"ADM": 0b111000,
	}
	a.compTable = map[string]uint16{
		// a == 0
		"0":   0b0_101010_000000,
		"1":   0b0_111111_000000,
		"-1":  0b0_111010_000000,
		"D":   0b0_001100_000000,
		"A":   0b0_110000_000000,
		"D+A": 0b0_000010_000000,
	}
	a.jumpTable = map[string]uint16{
		"JGT": 0b001,
		"JEQ": 0b010,
		"JGE": 0b011,
		"JLT": 0b100,
		"JNE": 0b101,
		"JLE": 0b110,
		"JMP": 0b111,
	}
	return a
}

func (a *Assembler) Compile(reader io.Reader) ([]byte, error) {
	buf := make([]byte, 0)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if !skip(line) {
			binary, err := a.compileLine(line)
			if err != nil {
				return nil, err
			}
			buf = append(buf, binary...)
		}
	}
	return buf, nil
}

func skip(line string) bool {
	s := strings.TrimSpace(line)
	if s == "" {
		return true
	}

	if strings.HasPrefix(s, "//") {
		return true
	}

	return false
}

func (a *Assembler) compile_a_instr(line string) ([]byte, error) {
	val, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, errors.New("invalid A instruction, not a valid number")
	}

	var ret uint16 = uint16(val & 0x7fff)
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, ret)
	return buf, nil
}

func (a *Assembler) compile_c_instr(line string) ([]byte, error) {
	var ret uint16 = 0xe000
	index0 := strings.Index(line, "=")
	index1 := strings.Index(line, ";")
	if index0 <= -1 {
		return nil, errors.New("invalid C instruction, '=' is missing")
	}

	dest := line[:index0]
	comp := ""
	jump := ""

	if index1 > -1 {
		jump = line[index1:]
		comp = line[index0+1 : index1-1]
	} else {
		comp = line[index0+1:]
	}

	value, exist := a.destTable[dest]
	if exist {
		ret |= value
	} else {
		return nil, errors.New("unknown dest")
	}

	if jump != "" {
		value, exist = a.jumpTable[jump]
		if exist {
			ret |= value
		} else {
			return nil, errors.New("unknown jump")
		}
	}

	value, exist = a.compTable[comp]
	if exist {
		ret |= value
	} else {
		return nil, errors.New("unknown comp")
	}

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, ret)
	return buf, nil
}

func (a *Assembler) compileLine(line string) ([]byte, error) {
	if line[0] == '@' {
		return a.compile_a_instr(line)
	} else {
		return a.compile_c_instr(line)
	}
}
