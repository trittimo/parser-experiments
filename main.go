package main

import (
	_ "embed"
	"fmt"
)

//go:embed tests/test.for
var parseText string

func main() {
	parser := CreateParser(&parseText)
	if tokens, err := parser.Parse(); err != nil {
		fmt.Printf("Error parsing text: %#v", err)
	} else {
		fmt.Printf("%+v", tokens)
	}
}
