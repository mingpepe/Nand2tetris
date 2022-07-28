package analyzer

import "io"

type VMWriter struct {
	writer io.Writer
}

func NewVMWriter(writer io.Writer) *VMWriter {
	vm := VMWriter{}
	vm.writer = writer
	return &vm
}

func (v *VMWriter) WritePush(segment, index int) {

}

func (v *VMWriter) WritePop(segment, index int) {

}

func (v *VMWriter) WriteArithmetic(cmd int) {

}

func (v *VMWriter) WriteLabel(lable string) {

}

func (v *VMWriter) WriteGoTo(lable string) {

}

func (v *VMWriter) WriteIf(label string) {

}

func (v *VMWriter) WriteCall(name string, nArgs int) {

}

func (v *VMWriter) WriteFunction(name string, nLocals int) {

}

func (v *VMWriter) WriteReturn() {

}
