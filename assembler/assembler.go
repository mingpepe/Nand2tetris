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
	builtInReg    map[string]uint16
	destTable     map[string]uint16
	compTable     map[string]uint16
	jumpTable     map[string]uint16
	labelTable    map[string]uint16
	symbolTable   map[string]uint16
	symbolAddress uint16
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
		"":    0b000_000,
		"M":   0b001_000,
		"D":   0b010_000,
		"DM":  0b011_000,
		"MD":  0b011_000,
		"A":   0b100_000,
		"AM":  0b101_000,
		"AD":  0b110_000,
		"ADM": 0b111_000,
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
		"":    0b000,
		"JGT": 0b001,
		"JEQ": 0b010,
		"JGE": 0b011,
		"JLT": 0b100,
		"JNE": 0b101,
		"JLE": 0b110,
		"JMP": 0b111,
	}
	a.labelTable = make(map[string]uint16)
	a.symbolTable = make(map[string]uint16)
	a.symbolAddress = 16
	return a
}

func (a *Assembler) getNewSymboAddress(variable string) uint16 {
	address := a.symbolAddress
	a.symbolTable[variable] = address
	a.symbolAddress++
	return address
}

func (a *Assembler) Compile(reader io.Reader) ([]byte, error) {
	buf := make([]byte, 0)
	scanner := bufio.NewScanner(reader)
	lines := make([]string, 0)
	var lineCount uint16 = 0
	for scanner.Scan() {
		line := scanner.Text()
		if !skip(line) {
			if strings.HasPrefix(line, "(") && strings.HasSuffix(line, ")") {
				label := line[1 : len(line)-1]
				a.labelTable[label] = lineCount
				fmt.Printf("Add label %s, value = %d\n", label, lineCount)
			} else {
				lines = append(lines, line)
			}
			lineCount += 1
		}
	}

	for _, line := range lines {
		binary, err := a.compileLine(line)
		if err != nil {
			return nil, err
		}
		buf = append(buf, binary...)
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
	variable := line[1:]
	val, exist := a.builtInReg[variable]
	if !exist {
		val, exist = a.labelTable[variable]
		if !exist {
			_val, err := strconv.Atoi(variable)
			if err == nil {
				val = uint16(_val)
			} else {
				val, exist = a.symbolTable[variable]
				if !exist {
					val = a.getNewSymboAddress(variable)
				}
			}
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

	dest := ""
	comp := ""
	jump := ""
	if index0 > -1 {
		dest = line[:index0]
		if index1 > -1 {
			jump = line[index1:]
			comp = line[:index1]
		} else {
			comp = line[index0+1:]
		}
	} else {
		if index1 > -1 {
			jump = line[index1+1:]
			comp = line[:index1]
		} else {
			comp = line
		}
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
