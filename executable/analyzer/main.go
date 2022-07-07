package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	out_filename := ""
	if *directory == "" {
		if !exist(*filename) {
			log.Printf("file not found: %s", *filename)
			return
		}

		if !strings.HasSuffix(*filename, ".jack") {
			log.Println("input must be a jack file")
			return
		}

		len := len(*filename)
		out_filename = (*filename)[:len-5] + "_KM.xml"
	} else {
		// Todo : handle multi files
		idx := strings.LastIndex(*directory, "\\")
		name := (*directory)[idx+1:]
		out_filename = *directory + "\\" + name + "xml"
	}
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Test only
	ana := analyzer.NewTokenAnalyzer(f)
	ana.Parse()

	print(out_filename + "\n")

	out_f, err := os.Create(out_filename)
	if err != nil {
		log.Print(err.Error())
	}
	defer out_f.Close()
	out_f.WriteString("<tokens>\n")
	for ana.HasMoreTokens() {
		token_type := ana.TokenType()
		token := ana.CurrentToken()
		_, err = out_f.WriteString(fmt.Sprintf("<%s>%s</%s>", token_type, token, token_type))
		if err != nil {
			log.Fatal(err.Error())
		}
		out_f.WriteString("\n")
		ana.Advance()
	}
	out_f.WriteString("</tokens>\n")
}
