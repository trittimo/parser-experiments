package main

import (
	_ "embed"
	"fmt"
	"strings"
)

type TokenForest struct {
	values []Token
}

func CreateTokenForest(tokens ...Token) TokenForest {
	return TokenForest{tokens}
}

func (t *TokenForest) Add(tokens ...Token) {
	t.values = append(t.values, tokens...)
}

func (t *TokenForest) Join(other TokenForest) int {
	t.values = append(t.values, other.values...)
	return len(other.values)
}

func (t *TokenForest) JoinIgnoreError(forest TokenForest, err error) bool {
	if err == nil {
		t.Join(forest)
		return true
	}
	return false
}

func (t TokenForest) AsString(tabCount int) string {
	switch len(t.values) {
	case 0:
		return "[]"
	default:
		{
			result := "[\n"
			for index, child := range t.values {
				result += child.AsString(tabCount + 1)
				if index < len(t.values)-1 {
					result += ",\n"
				} else {
					result += "\n"
				}
			}
			result += strings.Repeat("\t", tabCount) + "]"
			return result
		}
	}
}

func (t TokenForest) String() string {
	return t.AsString(0)
}

func (t *TokenForest) Len() int {
	return len(t.values)
}

// All tokens beyond this point
// ============================

type Token interface {
	AsString(tabCount int) string
}

type ProgramToken struct {
	name       string
	statements TokenForest
}

func Tab(tabCount int) string {
	return strings.Repeat("\t", tabCount)
}

func (t ProgramToken) AsString(tabCount int) string {
	tab := Tab(tabCount)
	switch t.statements.Len() {
	case 0:
		return fmt.Sprintf("%sPROGRAM(name: %s, statements: %s)", tab, t.name, t.statements.AsString(tabCount+1))
	default:
		return fmt.Sprintf("%sPROGRAM(name: %s, statements: %s\n%s)", tab, t.name, t.statements.AsString(tabCount+1), tab)
	}
}

type CommentToken struct {
	value string
}

func (t CommentToken) AsString(tabCount int) string {
	return fmt.Sprintf("%sCOMMENT(value: '%s')", Tab(tabCount), t.value)
}

//lint:ignore U1000 We may some day bring back whitespace token into the parse tree, so ignore this
type WhitespaceToken struct {
	value string
}

func (t WhitespaceToken) AsString(tabCount int) string {
	return Tab(tabCount) + "WHITESPACE()"
}

//lint:ignore U1000 We may some day bring back whitespace token into the parse tree, so ignore this
type NewLineToken struct {
	value string
}

func (t NewLineToken) AsString(tabCount int) string {
	return Tab(tabCount) + "NEWLINE()"
}

type ContinuationToken struct {
}

func (t ContinuationToken) AsString(tabCount int) string {
	return Tab(tabCount) + "CONTINUATION()"
}

type PrimitiveCallToken struct {
	kind   string
	values TokenForest
}

func (t PrimitiveCallToken) AsString(tabCount int) string {
	tab := Tab(tabCount)
	switch t.values.Len() {
	case 0:
		return fmt.Sprintf("%sPRIMITIVECALL(kind: %s, values: %s)", tab, t.kind, t.values.AsString(tabCount+1))
	default:
		return fmt.Sprintf("%sPRIMITIVECALL(kind: %s, values: %s\n%s)", tab, t.kind, t.values.AsString(tabCount+1), tab)
	}
}

type StringToken struct {
	value string
}

func (t StringToken) AsString(tabCount int) string {
	return Tab(tabCount) + fmt.Sprintf("STRING(value: '%s')", t.value)
}

type TypeToken struct {
	kind string
}

type FunctionToken struct {
	returnType  Token
	returnValue TokenForest
	statements  TokenForest
	name        string
}

func (t FunctionToken) AsString(tabCount int) string {
	tab := Tab(tabCount)
	statements := "[]"
	endparen := ")"
	returnType := t.returnType.AsString(tabCount + 1)
	returnValue := t.returnValue.AsString(tabCount + 1)
	if t.statements.Len() > 0 {
		statements = t.statements.AsString(tabCount + 1)
		endparen = tab + ")"
	}
	return tab + fmt.Sprintf("FUNCTION(name: %s, returnType: %s, returnValue: %s, statements: %s%s",
		t.name, returnType, returnValue, statements, endparen)
}
