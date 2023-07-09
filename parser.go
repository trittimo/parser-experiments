package main

import (
	"errors"
	"regexp"

	"unicode/utf8"
)

type ParseError struct {
	msg string
}

func (e *ParseError) Error() string {
	return e.msg
}

type Parser struct {
	cursors       []int
	currentCursor int
	text          []rune
}

func CreateParser(text *string) (parser Parser) {
	return Parser{make([]int, 16), 0, []rune(*text)}
}

func (p *Parser) Parse() (tokens TokenForest, err error) {
	child1 := ProgramToken{"child", CreateTokenForest()}
	child2 := ProgramToken{"child", CreateTokenForest()}
	child3 := ProgramToken{"child", CreateTokenForest(child1, child2)}
	tokens = CreateTokenForest(
		ProgramToken{"parent", CreateTokenForest(child1, child2)},
		ProgramToken{"parent", CreateTokenForest(child1, child2, child3)},
	)
	return
}

func (p *Parser) SaveCursor() {
	if p.currentCursor >= cap(p.cursors) {
		p.cursors = append(p.cursors, p.cursors[p.currentCursor])
	} else {
		p.cursors[p.currentCursor+1] = p.cursors[p.currentCursor]
	}
	p.currentCursor++
}
func (p *Parser) DeleteCursor() {
	p.currentCursor--
}
func (p *Parser) Cursor() (cursor int) {
	return p.cursors[p.currentCursor]
}
func (p *Parser) SetCursor(value int) {
	p.cursors[p.currentCursor] = value
}
func (p *Parser) DecrementCursor() {
	p.cursors[p.currentCursor]--
}
func (p *Parser) IncrementCursor() {
	p.cursors[p.currentCursor]++
}
func (p *Parser) CurrentCharacter() rune {
	return p.text[p.Cursor()]
}

func (p *Parser) CharacterAtIndex(index int) (r rune, err error) {
	if index >= len(p.text) || index < 0 {
		return -1, errors.New("")
	}
	return p.text[index], nil
}
func (p *Parser) Matches(other string) bool {
	var index int
	var otherChar rune
	for index, otherChar = range other {
		if c, err := p.CharacterAtIndex(p.Cursor() + index); err != nil {
			return false
		} else if otherChar != c {
			return false
		}
	}
	p.SetCursor(p.Cursor() + index + 1)
	return true
}

func (p *Parser) ReadRune() (r rune, size int, err error) {
	if r, err = p.CharacterAtIndex(p.Cursor()); err != nil {
		return -1, -1, err
	}
	size = utf8.RuneLen(r)
	p.IncrementCursor()
	return
}

func (p *Parser) MatchRange(loc int[]) string {
	return string(p.text(p.Cursor()+loc[0] : p.Cursor()+loc[1]))
}

var literalRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_$]*`)

func (p *Parser) Literal() (literal string, err error) {
	p.SaveCursor()
	if matchResult := literalRegex.FindReaderIndex(p); matchResult == nil {
		err = errors.New("expected literal")
		p.DeleteCursor()
	} else {
		literal = Range(matchResult)
	}
	return
}

// TODO: split this up so the first part matches the newline, then we match for n comments, then match continuation
var continuationTokenRegex = regexp.MustCompile(`^([\r\n]+(\s*!.*)?)+\s*1\s+`)
func (p *Parser) AcceptContinuationToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	if matchResult := continuationTokenRegex.FindReaderIndex(p); matchResult == nil {
		acceptError = errors.New("expected continuation")
		p.DeleteCursor()
	} else {
		forest.Add(ContinuationToken{value: p.Range(matchResult)})
	}
	return
}

var whitespaceRegex = regexp.MustCompile(`^[ \t]`)
func (p *Parser) AcceptWhitespaceToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	if matchResult := whitespaceRegex.FindReaderIndex(p); matchResult == nil {
		if continuationForest, err := p.AcceptContinuationToken(); err != nil {
			forest.Join(continuationForest)
			return
		}
		acceptError = errors.New("expected whitespace")
		p.DeleteCursor()
	} else {
		forest.Add(WhitespaceToken{value: p.Range(matchResult)})
		forest.JoinIgnoreError(p.AcceptContinuationToken())
	}
	return
}

// TODO: Add any discovered comments/whitespace to the forest here
var newLineRegex = regexp.MustCompile(`^[\r\n]+`)
func (p *Parser) AcceptNewLineToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	if matchResult := newLineRegex.FindReaderIndex(p); matchResult == nil {
		acceptError = errors.New("expected newline(s)")
		p.DeleteCursor()
	} else {
		forest.Add(NewLineToken{value: p.Range(matchResult)})
	}
	return
}


var commentRegex = regexp.MustCompile(`^!.*`)
func (p *Parser) AcceptCommentToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	forest.JoinIgnoreError(p.AcceptWhitespaceToken())

	if matchResult := commentRegex.FindReaderIndex(p); matchResult == nil {
		acceptError = errors.New("expected comment")
		p.DeleteCursor()
	} else {
		forest.Add(CommentToken{})
	}
}

func (p *Parser) AcceptStatementToken() (forest TokenForest, acceptError error) {
	
}

func (p *Parser) AcceptProgramToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	if !p.Matches("program") {
		p.DeleteCursor()
		acceptError = errors.New("expected 'program'")
		return
	}
	if name, err := p.Literal(); err != nil {
		p.DeleteCursor()
		acceptError = err
		return
	} else {
		token := ProgramToken{name: name, statements: CreateTokenForest()}
		for {
			if statement, err := p.AcceptStatementToken(); err != nil {
				break
			} else {
				token.statements.Join(statement)
			}
		}
		forest.Add(token)
		return
	}
}
