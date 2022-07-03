package vm

import (
	"bufio"
	"io"
	"strings"
)

type VM struct {
}

func New() *VM {
	return nil
}

func (vm *VM) Compile(reader io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(reader)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if !skip(line) {
			lines = append(lines, line)
		}
	}
	return nil, nil
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
