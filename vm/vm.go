package vm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

const (
	C_ARITHMETIC int = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

const arithmeticTemplate = "@SP\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"A=A-1\n"

type VM struct {
	arthJumpFlag    int
	retLabelCnt     int
	currentFilename string
}

func New() *VM {
	vm := new(VM)
	vm.arthJumpFlag = 0
	vm.retLabelCnt = 0
	vm.currentFilename = ""
	return vm
}

func (vm *VM) BootstrapCode() string {
	tmp := "@256\n" +
		"D=A\n" +
		"@SP\n" +
		"M=D\n"
	return tmp + vm.compile_line("call Sys.init 0")
}

func (vm *VM) Compile(filename string, reader io.Reader) (string, error) {
	vm.currentFilename = filename
	scanner := bufio.NewScanner(reader)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		// Remove comments
		idx := strings.Index(line, "//")
		if idx > 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if !skip(line) {
			lines = append(lines, line)
		}
	}

	asm := ""
	for i := 0; i < len(lines); i++ {
		asm += "//" + lines[i] + "\n"
		asm += vm.compile_line(lines[i])
	}
	return asm, nil
}

func (vm *VM) compile_line(line string) string {
	cmd_type, err := getCmdType(line)
	if err != nil {
		log.Fatal(err)
	}
	switch cmd_type {
	case C_ARITHMETIC:
		cmd := strings.Split(line, " ")[0]
		switch cmd {
		case "add":
			return arithmeticTemplate + "M=M+D\n"
		case "sub":
			return arithmeticTemplate + "M=M-D\n"
		case "and":
			return arithmeticTemplate + "M=M&D\n"
		case "or":
			return arithmeticTemplate + "M=M|D\n"
		case "not":
			return "@SP\nA=M-1\nM=!M\n"
		case "neg":
			return "D=0\n@SP\nA=M-1\nM=D-M\n"
		case "gt":
			{
				asm := generateArithCompareCode("JLE", vm.arthJumpFlag)
				vm.arthJumpFlag++
				return asm
			}
		case "lt":
			{
				asm := generateArithCompareCode("JGE", vm.arthJumpFlag)
				vm.arthJumpFlag++
				return asm
			}
		case "eq":
			{
				asm := generateArithCompareCode("JNE", vm.arthJumpFlag)
				vm.arthJumpFlag++
				return asm
			}
		}
	case C_PUSH:
		{
			segment, err := getArg1(line)
			if err != nil {
				log.Fatal(err)
			}
			idx, err := getArg2(line)
			if err != nil {
				log.Fatal(err)
			}
			switch segment {
			case "constant":
				return fmt.Sprintf("@%d\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", idx)
			case "local":
				return generatePointerPushCode("LCL", idx)
			case "argument":
				return generatePointerPushCode("ARG", idx)
			case "this":
				return generatePointerPushCode("THIS", idx)
			case "that":
				return generatePointerPushCode("THAT", idx)
			case "temp":
				return generateDirectPushCode(fmt.Sprintf("%d", idx+5))
			case "pointer":
				if idx == 0 {
					return generateDirectPushCode("THIS")
				} else if idx == 1 {
					return generateDirectPushCode("THAT")
				}
			case "static":
				return generateDirectPushCode(fmt.Sprintf("%s.%d", vm.currentFilename, idx))
			}
		}
	case C_POP:
		segment, err := getArg1(line)
		if err != nil {
			log.Fatal(err)
		}
		idx, err := getArg2(line)
		if err != nil {
			log.Fatal(err)
		}
		switch segment {
		case "local":
			return generatePointerPopCode("LCL", idx)
		case "argument":
			return generatePointerPopCode("ARG", idx)
		case "this":
			return generatePointerPopCode("THIS", idx)
		case "that":
			return generatePointerPopCode("THAT", idx)
		case "temp":
			return generateDirectPopCode(fmt.Sprintf("%d", idx+5))
		case "pointer":
			if idx == 0 {
				return generateDirectPopCode("THIS")
			} else if idx == 1 {
				return generateDirectPopCode("THAT")
			}
		case "static":
			return generateDirectPopCode(fmt.Sprintf("%s.%d", vm.currentFilename, idx))
		}
	case C_LABEL:
		{
			name, err := getArg1(line)
			if err != nil {
				log.Fatal(err)
			}
			return fmt.Sprintf("(%s)\n", name)
		}
	case C_GOTO:
		{
			name, err := getArg1(line)
			if err != nil {
				log.Fatal(err)
			}
			return fmt.Sprintf("@%s\n0;JMP\n", name)
		}
	case C_IF:
		{
			name, err := getArg1(line)
			if err != nil {
				log.Fatal(err)
			}
			return fmt.Sprintf("%s@%s\nD;JNE\n", arithmeticTemplate, name)
		}
	case C_FUNCTION:
		{
			name, err := getArg1(line)
			if err != nil {
				log.Fatal(err)
			}
			numArgs, err := getArg2(line)
			if err != nil {
				log.Fatal(err)
			}
			tmp := fmt.Sprintf("(%s)\n", name)
			for i := 0; i < numArgs; i++ {
				// The same with push constant 0
				tmp += fmt.Sprintf("@%d\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", 0)
			}
			return tmp
		}
	case C_RETURN:
		return generateReturnCode()
	case C_CALL:
		{
			name, err := getArg1(line)
			if err != nil {
				log.Fatal(err)
			}
			numArgs, err := getArg2(line)
			if err != nil {
				log.Fatal(err)
			}
			newLabel := fmt.Sprintf("RETURN_LABEL%d", vm.retLabelCnt)
			vm.retLabelCnt++
			return fmt.Sprintf("@%s\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", newLabel) + // push return address
				generateDirectPushCode("LCL") +
				generateDirectPushCode("ARG") +
				generateDirectPushCode("THIS") +
				generateDirectPushCode("THAT") +
				"@SP\n" +
				"D=M\n" +
				"@5\n" +
				"D=D-A\n" +
				"@" + strconv.Itoa(numArgs) + "\n" +
				"D=D-A\n" +
				"@ARG\n" +
				"M=D\n" +
				"@SP\n" +
				"D=M\n" +
				"@LCL\n" +
				"M=D\n" +
				"@" + name + "\n" +
				"0;JMP\n" +
				"(" + newLabel + ")\n"
		}
	}
	return "unexpected return : " + line + "\n"
}

func generateArithCompareCode(_type string, arthJumpFlag int) string {
	return "@SP\n" +
		"AM=M-1\n" +
		"D=M\n" +
		"A=A-1\n" +
		"D=M-D\n" +
		"@FALSE" + strconv.Itoa(arthJumpFlag) + "\n" +
		"D;" + _type + "\n" +
		"@SP\n" +
		"A=M-1\n" +
		"M=-1\n" +
		"@CONTINUE" + strconv.Itoa(arthJumpFlag) + "\n" +
		"0;JMP\n" +
		"(FALSE" + strconv.Itoa(arthJumpFlag) + ")\n" +
		"@SP\n" +
		"A=M-1\n" +
		"M=0\n" +
		"(CONTINUE" + strconv.Itoa(arthJumpFlag) + ")\n"
}

func generatePointerPushCode(seg string, idx int) string {
	return fmt.Sprintf("@%s\n", seg) +
		"D=M\n" +
		fmt.Sprintf("@%d\nA=D+A\nD=M\n", idx) +
		"@SP\n" +
		"A=M\n" +
		"M=D\n" +
		"@SP\n" +
		"M=M+1\n"
}

func generateDirectPushCode(seg string) string {
	return fmt.Sprintf("@%s\n", seg) +
		"D=M\n" +
		"@SP\n" +
		"A=M\n" +
		"M=D\n" +
		"@SP\n" +
		"M=M+1\n"
}

func generatePointerPopCode(seg string, idx int) string {
	return fmt.Sprintf("@%s\nD=M\n@%d\nD=D+A\n", seg, idx) +
		"@R13\n" +
		"M=D\n" +
		"@SP\n" +
		"AM=M-1\n" +
		"D=M\n" +
		"@R13\n" +
		"A=M\n" +
		"M=D\n"
}

func generateDirectPopCode(seg string) string {
	return fmt.Sprintf("@%s\nD=A\n", seg) +
		"@R13\n" +
		"M=D\n" +
		"@SP\n" +
		"AM=M-1\n" +
		"D=M\n" +
		"@R13\n" +
		"A=M\n" +
		"M=D\n"
}

func getCmdType(line string) (int, error) {
	type_string := make([]string, 0)
	type_val := make([]int, 0)
	// Arithmetic
	type_string = append(type_string, "add")
	type_val = append(type_val, C_ARITHMETIC)
	type_string = append(type_string, "sub")
	type_val = append(type_val, C_ARITHMETIC)
	type_string = append(type_string, "neg")
	type_val = append(type_val, C_ARITHMETIC)
	type_string = append(type_string, "eq")
	type_val = append(type_val, C_ARITHMETIC)
	type_string = append(type_string, "gt")
	type_val = append(type_val, C_ARITHMETIC)
	type_string = append(type_string, "lt")
	type_val = append(type_val, C_ARITHMETIC)
	type_string = append(type_string, "and")
	type_val = append(type_val, C_ARITHMETIC)
	type_string = append(type_string, "or")
	type_val = append(type_val, C_ARITHMETIC)
	type_string = append(type_string, "not")
	type_val = append(type_val, C_ARITHMETIC)
	// Others
	type_string = append(type_string, "push")
	type_val = append(type_val, C_PUSH)
	type_string = append(type_string, "pop")
	type_val = append(type_val, C_POP)
	type_string = append(type_string, "label")
	type_val = append(type_val, C_LABEL)
	type_string = append(type_string, "goto")
	type_val = append(type_val, C_GOTO)
	type_string = append(type_string, "if-goto")
	type_val = append(type_val, C_IF)
	type_string = append(type_string, "func")
	type_val = append(type_val, C_FUNCTION)
	type_string = append(type_string, "return")
	type_val = append(type_val, C_RETURN)
	type_string = append(type_string, "call")
	type_val = append(type_val, C_CALL)

	for i := 0; i < len(type_string); i++ {
		if strings.HasPrefix(line, type_string[i]) {
			return type_val[i], nil
		}
	}

	return -1, errors.New("unknown cmd type : " + line)
}

func getArg1(line string) (string, error) {
	// Should not call for C_RETURN
	sep := strings.Split(line, " ")
	if len(sep) < 2 {
		return "", errors.New("unexpected data(get_arg1) : " + line)
	}
	return sep[1], nil
}

func getArg2(line string) (int, error) {
	// Only for C_PUSH, C_POP, C_FUNCTION, C_CALL
	sep := strings.Split(line, " ")
	if len(sep) != 3 {
		return 0, errors.New("unexpected data(get_arg2) : " + line)
	}
	val, err := strconv.Atoi(sep[2])
	if err != nil {
		return -1, errors.New("unexpected data(get_arg2) : " + line)
	}
	return val, nil
}

func skip(line string) bool {
	if line == "" {
		return true
	}
	if strings.HasPrefix(line, "//") {
		return true
	}

	return false
}

func preFrameTemplate(position string) string {
	return "@R11\n" +
		"D=M-1\n" +
		"AM=D\n" +
		"D=M\n" +
		"@" + position + "\n" +
		"M=D\n"
}

func generateReturnCode() string {
	return "@LCL\n" +
		"D=M\n" +
		"@R11\n" +
		"M=D\n" +
		"@5\n" +
		"A=D-A\n" +
		"D=M\n" +
		"@R12\n" +
		"M=D\n" +
		generatePointerPopCode("ARG", 0) +
		"@ARG\n" +
		"D=M\n" +
		"@SP\n" +
		"M=D+1\n" +
		preFrameTemplate("THAT") +
		preFrameTemplate("THIS") +
		preFrameTemplate("ARG") +
		preFrameTemplate("LCL") +
		"@R12\n" +
		"A=M\n" +
		"0;JMP\n"
}
