package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mingpepe/Nand2teris/analyzer"
)

func exist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func main() {
	var filename = flag.String("f", "input.jack", "input filename")
	var directory = flag.String("d", "", "directory contains jack files")
	flag.Parse()

	filenames := make([]string, 0)
	if *directory == "" {
		if !exist(*filename) {
			log.Printf("file not found: %s", *filename)
			return
		}

		if !strings.HasSuffix(*filename, ".jack") {
			log.Println("input must be a jack file")
			return
		}
		filenames = append(filenames, *filename)
	} else {
		err := filepath.Walk(*directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if strings.HasSuffix(path, ".jack") {
					filenames = append(filenames, path)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, _filename := range filenames {
		f, err := os.Open(_filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		ana := analyzer.NewTokenAnalyzer(f)
		ana.Parse()

		length := len(_filename)
		out_filename := (_filename)[:length-5] + "_KM.xml"

		out_f, err := os.Create(out_filename)
		if err != nil {
			log.Print(err.Error())
		}
		defer out_f.Close()
	}
}
