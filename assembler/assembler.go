package assembler

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Assembler struct {
	builtInReg map[string]uint16
	destTable  map[string]uint16
	compTable  map[string]uint16
	jumpTable  map[string]uint16
}

func New() *Assembler {
	a := &Assembler{}
	a.builtInReg = map[string]uint16{
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15,
		"SCREEN": 16384,
		"KBD":    24576,
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
	}
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
		// a = 0
		"0":   0b0_101010_000000,
		"1":   0b0_111111_000000,
		"-1":  0b0_111010_000000,
		"D":   0b0_001100_000000,
		"A":   0b0_110000_000000,
		"-D":  0b0_001101_000000,
		"-A":  0b0_110001_000000,
		"D+1": 0b0_011111_000000,
		"A+1": 0b0_110111_000000,
		"D-1": 0b0_001110_000000,
		"A-1": 0b0_110010_000000,
		"D+A": 0b0_000010_000000,
		"D-A": 0b0_010011_000000,
		"A-D": 0b0_000111_000000,
		"D&A": 0b0_000000_000000,
		"D|A": 0b0_010101_000000,
		// a = 1
		"M":   0b1_110000_000000,
		"!M":  0b1_110001_000000,
		"-M":  0b1_110011_000000,
		"M+1": 0b1_110111_000000,
		"M-1": 0b1_110010_000000,
		"D+M": 0b1_000010_000000,
		"D-M": 0b1_010011_000000,
		"M-D": 0b1_000111_000000,
		"D&M": 0b1_000000_000000,
		"D|M": 0b1_010101_000000,
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
	val, exist := a.builtInReg[line[1:]]
	if !exist {
		_val, err := strconv.Atoi(line[1:])
		if err == nil {
			val = uint16(_val)
		} else {
			return nil, fmt.Errorf("invalid A instruction, not a valid number(%s), nor build-in register", line)
		}
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
		return nil, fmt.Errorf("invalid C instruction, '=' is missing(%s)", line)
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
		return nil, fmt.Errorf("unknown dest(dest = %s, line = %s)", dest, line)
	}

	if jump != "" {
		value, exist = a.jumpTable[jump]
		if exist {
			ret |= value
		} else {
			return nil, fmt.Errorf("unknown jump(jump = %s, line = %s)", jump, line)
		}
	}

	value, exist = a.compTable[comp]
	if exist {
		ret |= value
	} else {
		return nil, fmt.Errorf("unknown comp(comp = %s, line = %s)", comp, line)
	}

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, ret)
	return buf, nil
}

func (a *Assembler) compileLine(line string) ([]byte, error) {
	index := strings.Index(line, "//")
	if index != -1 {
		line = line[:index-1]
	}
	line = strings.TrimSpace(line)
	if line[0] == '@' {
		return a.compile_a_instr(line)
	} else {
		return a.compile_c_instr(line)
	}
}
