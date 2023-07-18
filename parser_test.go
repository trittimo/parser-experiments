package main

import "testing"

func TestWhitespace(t *testing.T) {
	text := "  test"
	parser := CreateParser(&text)
	if forest, err := parser.AcceptWhitespaceToken(); err != nil {
		t.Fatalf("Whitespace not accepted: %s", err)
	} else if forest.Len() != 0 {
		t.Fatalf("Unexpected forest length: %d", forest.Len())
	}
	if parser.Cursor() != 2 {
		t.Fatalf("Unexpected cursor position: %d", parser.Cursor())
	}
}

func TestString(t *testing.T) {
	text := "'this is a string'"
	compare := StringToken{value: "this is a string"}
	parser := CreateParser(&text)
	if forest, err := parser.AcceptStringToken(); err != nil {
		t.Fatalf("String not accepted: %s", err)
	} else if forest.Len() != 1 {
		t.Fatalf("Unexpected forest length: %d", forest.Len())
	} else if forest.values[0].AsString(0) != compare.AsString(0) {
		t.Fatalf("String did not match: expected %s got %s", compare.AsString(0), forest.values[0].AsString(0))
	}
}
