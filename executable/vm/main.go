package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
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
	var directory = flag.String("d", "", "directory contains vm files")
	flag.Parse()

	filenames := make([]string, 0)
	out_filename := ""
	if *directory == "" {
		if !exist(*filename) {
			log.Printf("file not found: %s", *filename)
			return
		}

		if !strings.HasSuffix(*filename, ".vm") {
			log.Printf("input must be a vm file")
			return
		}
		filenames = append(filenames, *filename)

	} else {
		err := filepath.Walk(*directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if strings.HasSuffix(path, ".vm") {
					filenames = append(filenames, path)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	v := vm.New()
	code := ""
	if !*bypass_bootstrap {
		code += v.BootstrapCode()
	}

	for _, _filename := range filenames {
		f, err := os.Open(_filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		asm, err := v.Compile(f)
		if err != nil {
			log.Print(err.Error())
		}
		code += asm
	}

	if *directory == "" {
		idx := strings.LastIndex(*filename, ".")
		out_filename = (*filename)[:idx] + ".asm"
	} else {
		idx := strings.LastIndex(*directory, "\\")
		dir_name := (*directory)[idx+1:]
		out_filename = *directory + "\\" + dir_name + ".asm"
	}
	out_f, err := os.Create(out_filename)
	if err != nil {
		log.Print(err.Error())
	}
	defer out_f.Close()
	_, err = out_f.WriteString(code)
	if err != nil {
		log.Print(err.Error())
	}
}
