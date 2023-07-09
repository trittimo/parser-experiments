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

func (t *TokenForest) Join(other TokenForest) {
	t.values = append(t.values, other.values...)
}

func (t *TokenForest) JoinIgnoreError(forest TokenForest, err error) {
	if err == nil {
		t.Join(forest)
	}
}

func (t TokenForest) StringTabbed(tabCount int) string {
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
	return t.StringTabbed(0)
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
		return fmt.Sprintf("%sPROGRAM(name: %s, statements: %s)", tab, t.name, t.statements.StringTabbed(tabCount+1))
	default:
		return fmt.Sprintf("%sPROGRAM(name: %s, statements: %s\n%s)", tab, t.name, t.statements.StringTabbed(tabCount+1), tab)
	}
}

type CommentToken struct {
	value string
}

func (t CommentToken) AsString(tabCount int) string {
	return Tab(tabCount) + "COMMENT(value: '" + t.value + "')"
}

type WhitespaceToken struct {
	value string
}

func (t WhitespaceToken) AsString(tabCount int) string {
	return Tab(tabCount) + "WHITESPACE()"
}

type NewLineToken struct {
	value string
}

func (t NewLineToken) AsString(tabCount int) string {
	return Tab(tabCount) + "NEWLINE()"
}

type ContinuationToken struct {
	value string
}

func (t ContinuationToken) AsString(tabCount int) string {
	return Tab(tabCount) + "CONTINUATION()"
}
