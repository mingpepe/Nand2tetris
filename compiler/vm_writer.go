package compiler

import (
	"io"
	"strconv"
)

type VMWriter struct {
	writer io.Writer
}

func NewVMWriter(writer io.Writer) *VMWriter {
	vm := VMWriter{}
	vm.writer = writer
	return &vm
}

func (v *VMWriter) WritePush(segment string, index int) {
	v.writeVMCode("push", segment, strconv.Itoa(index))
}

func (v *VMWriter) WritePop(segment string, index int) {
	v.writeVMCode("pop", segment, strconv.Itoa(index))
}

func (v *VMWriter) WriteArithmetic(cmd string) {
	v.writeVMCode(cmd)
}

func (v *VMWriter) WriteLabel(lable string) {
	v.writeVMCode("label", lable)
}

func (v *VMWriter) WriteGoTo(lable string) {
	v.writeVMCode("goto", lable)
}

func (v *VMWriter) WriteIf(lable string) {
	v.writeVMCode("if-goto", lable)
}

func (v *VMWriter) WriteCall(name string, nArgs int) {
	v.writeVMCode("call", name, strconv.Itoa(nArgs))
}

func (v *VMWriter) WriteFunction(name string, nLocals int) {
	v.writeVMCode("function", name, strconv.Itoa(nLocals))
}

func (v *VMWriter) WriteReturn() {
	v.writeVMCode("return")
}

func (v *VMWriter) WriteComment(comment string) {
	v.writer.Write([]byte("// " + comment + "\n"))
}

func (v *VMWriter) writeVMCode(cmd string, args ...string) {
	if len(args) >= 1 {
		cmd += " " + args[0]
	}
	if len(args) >= 2 {
		cmd += " " + args[1]
	}
	cmd += "\n"
	v.writer.Write([]byte(cmd))
}
