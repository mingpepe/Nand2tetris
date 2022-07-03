package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mingpepe/Nand2teris/vm"
)

func exist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func main() {
	var filename = flag.String("f", "input.vm", "input filename")
	flag.Parse()

	if !exist(*filename) {
		log.Printf("file not found: %s", *filename)
		return
	}

	if !strings.HasSuffix(*filename, ".vm") {
		log.Printf("input must be a vm file")
		return
	}

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	v := vm.New()
	asm, err := v.Compile(f)
	if err != nil {
		log.Print(err.Error())
	}

	fmt.Println(asm)
}
