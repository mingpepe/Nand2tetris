package main

import (
	"flag"
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
	var bypass_bootstrap = flag.Bool("bypass", false, "bypass bootstrap code for test")
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
	boolstrap := ""
	if !*bypass_bootstrap {
		boolstrap = v.BootstrapCode()
	}
	asm, err := v.Compile(f)
	if err != nil {
		log.Print(err.Error())
	}

	idx := strings.LastIndex(*filename, ".")
	out_filename := (*filename)[:idx] + ".asm"

	out_f, err := os.Create(out_filename)
	if err != nil {
		log.Print(err.Error())
	}
	defer out_f.Close()
	_, err = out_f.WriteString(boolstrap + asm)
	if err != nil {
		log.Print(err.Error())
	}
	f.Sync()
}
