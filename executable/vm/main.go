package main

import (
	"log"

	"github.com/mingpepe/Nand2teris/vm"
)

func main() {
	v := vm.New()
	_, err := v.Compile(nil)
	if err != nil {
		log.Print(err.Error())
	}
}
