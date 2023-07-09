package main

import (
	_ "embed"
	"fmt"
)

//go:embed blob/test.for
var parseText string

func main() {
	parser := CreateParser(&parseText)
	if tokens, err := parser.Parse(); err != nil {
		fmt.Printf("Error parsing text: %s", err.Error())
	} else {
		fmt.Printf("%+v", tokens)
	}
}
