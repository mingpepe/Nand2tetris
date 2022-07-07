package main

import (
	"flag"
	"log"
	"os"
	"strings"
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
		out_filename = (*filename)[:len-4] + ".xml"
	} else {
		idx := strings.LastIndex(*directory, "\\")
		name := (*directory)[idx+1:]
		out_filename = *directory + "\\" + name + ".xml"
	}
	print(out_filename + "\n")

	out_f, err := os.Create(out_filename)
	if err != nil {
		log.Print(err.Error())
	}
	defer out_f.Close()
	_, err = out_f.WriteString("test")
	if err != nil {
		log.Print(err.Error())
	}
}
