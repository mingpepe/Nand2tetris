package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	dir := filepath.Dir(*filename)
	base := filepath.Base(*filename)
	name := base[:len(base)-len(filepath.Ext(base))]

	newExt := ".hack"
	outputPath := filepath.Join(dir, name+newExt)

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if len(binary)%2 != 0 {
		log.Printf("Length of binary should be even")
		return
	}

	for i := 0; i < len(binary); i += 2 {
		fmt.Fprintf(file, "%08b%08b\n", binary[i], binary[i+1])
	}
}
