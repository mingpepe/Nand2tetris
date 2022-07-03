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

const arithmeticTemplate1 = "@SP\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"A=A-1\n"

type VM struct {
	arthJumpFlag  int
	ret_label_cnt int
}

func New() *VM {
	vm := new(VM)
	vm.arthJumpFlag = 0
	vm.ret_label_cnt = 0
	return vm
}

func (vm *VM) Compile(reader io.Reader) (string, error) {
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
		asm += vm.compile_line(lines[i])
	}
	return asm, nil
}

func (vm *VM) compile_line(line string) string {
	type1, err := get_cmd_type(line)
	if err != nil {
		log.Fatal(err)
	}
	switch type1 {
	case C_ARITHMETIC:
		cmd := strings.Split(line, " ")[0]
		switch cmd {
		case "add":
			return arithmeticTemplate1 + "M=M+D\n"
		case "sub":
			return arithmeticTemplate1 + "M=M-D\n"
		case "and":
			return arithmeticTemplate1 + "M=M&D\n"
		case "or":
			return arithmeticTemplate1 + "M=M|D\n"
		case "not":
			return "@SP\nA=M-1\nM=!M\n"
		case "neg":
			return "D=0\n@SP\nA=M-1\nM=D-M\n"
		case "gt":
			{
				asm := vm.arithmeticTemplate2("JLE")
				vm.arthJumpFlag++
				return asm
			}
		case "lt":
			{
				asm := vm.arithmeticTemplate2("JGE")
				vm.arthJumpFlag++
				return asm
			}
		case "eq":
			{
				asm := vm.arithmeticTemplate2("JNE")
				vm.arthJumpFlag++
				return asm
			}
		}
	case C_PUSH:
		{
			segment, err := get_arg1(line)
			if err != nil {
				log.Fatal(err)
			}
			idx, err := get_arg2(line)
			if err != nil {
				log.Fatal(err)
			}
			switch segment {
			case "constant":
				return fmt.Sprintf("@%d\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", idx)
			case "local":
				return vm.push_template("LCL", idx, true)
			case "argument":
				return vm.push_template("ARG", idx, true)
			case "this":
				return vm.push_template("THIS", idx, true)
			case "that":
				return vm.push_template("THAT", idx, true)
			case "temp":
				return vm.push_template("R5", idx+5, true)
			case "pointer":
				if idx == 0 {
					return vm.push_template("THIS", idx, false)
				} else if idx == 1 {
					return vm.push_template("THAT", idx, false)
				}
			case "static":
				return vm.push_template(fmt.Sprintf("%d", idx+16), idx, false)
			}
		}
	case C_POP:
		segment, err := get_arg1(line)
		if err != nil {
			log.Fatal(err)
		}
		idx, err := get_arg2(line)
		if err != nil {
			log.Fatal(err)
		}
		switch segment {
		case "local":
			return pop_template("LCL", idx, true)
		case "argument":
			return pop_template("ARG", idx, true)
		case "this":
			return pop_template("THIS", idx, true)
		case "that":
			return pop_template("THAT", idx, true)
		case "temp":
			return pop_template("R5", idx+5, true)
		case "pointer":
			if idx == 0 {
				return pop_template("THIS", idx, false)
			} else if idx == 1 {
				return pop_template("THAT", idx, false)
			}
		case "static":
			return pop_template(fmt.Sprintf("%d", idx+16), idx, false)
		}
	case C_LABEL:
		{
			name, err := get_arg1(line)
			if err != nil {
				log.Fatal(err)
			}
			return fmt.Sprintf("(%s)", name)
		}
	case C_GOTO:
		{
			name, err := get_arg1(line)
			if err != nil {
				log.Fatal(err)
			}
			return fmt.Sprintf("@%s\n0;JNE\n", name)
		}
	case C_IF:
		{
			name, err := get_arg1(line)
			if err != nil {
				log.Fatal(err)
			}
			return fmt.Sprintf("%s@%s\nD;JNE\n", arithmeticTemplate1, name)
		}
	case C_FUNCTION:
		{
			name, err := get_arg1(line)
			if err != nil {
				log.Fatal(err)
			}
			numArgs, err := get_arg2(line)
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
		return returnTemplate()
	case C_CALL:
		{
			name, err := get_arg1(line)
			if err != nil {
				log.Fatal(err)
			}
			numArgs, err := get_arg2(line)
			if err != nil {
				log.Fatal(err)
			}
			newLabel := fmt.Sprintf("RETURN_LABEL%d", vm.ret_label_cnt)
			vm.ret_label_cnt++
			return fmt.Sprintf("@%s\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", newLabel) + // push return address
				vm.push_template("LCL", 0, false) +
				vm.push_template("ARG", 0, false) +
				vm.push_template("THIS", 0, false) +
				vm.push_template("THAT", 0, false) +
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

func (vm *VM) arithmeticTemplate2(_type string) string {
	return "@SP\n" +
		"AM=M-1\n" +
		"D=M\n" +
		"A=A-1\n" +
		"D=M-D\n" +
		"@FALSE" + strconv.Itoa(vm.arthJumpFlag) + "\n" +
		"D;" + _type + "\n" +
		"@SP\n" +
		"A=M-1\n" +
		"M=-1\n" +
		"@CONTINUE" + strconv.Itoa(vm.arthJumpFlag) + "\n" +
		"0;JMP\n" +
		"(FALSE" + strconv.Itoa(vm.arthJumpFlag) + ")\n" +
		"@SP\n" +
		"A=M-1\n" +
		"M=0\n" +
		"(CONTINUE" + strconv.Itoa(vm.arthJumpFlag) + ")\n"
}

func (vm *VM) push_template(seg string, idx int, is_pointer bool) string {
	pointer_code := ""
	if is_pointer {
		pointer_code = fmt.Sprintf("@%d\nA=D+A\nD=M\n", idx)
	}
	return fmt.Sprintf("@%s\n", seg) +
		"D=M\n" +
		pointer_code +
		"@SP\n" +
		"A=M\n" +
		"M=D\n" +
		"@SP\n" +
		"M=M+1\n"
}

func pop_template(seg string, idx int, is_pointer bool) string {
	pointer_code := "D=A\n"
	if is_pointer {
		pointer_code = fmt.Sprintf("D=M\n@%d\nD=D+A\n", idx)
	}
	return fmt.Sprintf("@%s\n", seg) +
		pointer_code +
		"@R13\n" +
		"M=D\n" +
		"@SP\n" +
		"AM=M-1\n" +
		"D=M\n" +
		"@R13\n" +
		"A=M\n" +
		"M=D\n"
}

func get_cmd_type(line string) (int, error) {
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

func get_arg1(line string) (string, error) {
	// Should not call for C_RETURN
	sep := strings.Split(line, " ")
	if len(sep) < 2 {
		return "", errors.New("unexpected data(get_arg1) : " + line)
	}
	return sep[1], nil
}

func get_arg2(line string) (int, error) {
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

func returnTemplate() string {
	return "@LCL\n" +
		"D=M\n" +
		"@R11\n" +
		"M=D\n" +
		"@5\n" +
		"A=D-A\n" +
		"D=M\n" +
		"@R12\n" +
		"M=D\n" +
		pop_template("ARG", 0, true) +
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
