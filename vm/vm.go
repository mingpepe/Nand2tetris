package vm

import (
	"io"
)

type VM struct {
}

func New() *VM {
	return nil
}

func (vm *VM) Compile(reader io.Reader) ([]byte, error) {
	return nil, nil
}
