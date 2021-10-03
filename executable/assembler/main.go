package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mingpepe/Nand2teris/assembler"
)

func exist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func main() {
	var filename = flag.String("f", "input.asm", "input filename")
	flag.Parse()

	if !exist(*filename) {
		log.Printf("file not found: %s", *filename)
		return
	}

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	assemb := assembler.New()
	binary, err := assemb.Compile(f)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(binary); i += 2 {
		fmt.Printf("%08b%08b\n", binary[i], binary[i+1])
	}
}
