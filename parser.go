package main

import (
	"errors"
	"regexp"
	"strings"

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

func (p *Parser) Parse() (forest TokenForest, acceptError error) {
	if f, err := p.AcceptStatementTokens(); err != nil {
		acceptError = err
	} else if f.Len() == 0 {
		acceptError = errors.New("no statements found")
	} else {
		forest = f
	}
	return
}

func (p *Parser) SaveCursor() int {
	if p.currentCursor+1 >= cap(p.cursors) {
		p.cursors = append(p.cursors, p.cursors[p.currentCursor])
		p.cursors = p.cursors[:cap(p.cursors)]
	} else {
		p.cursors[p.currentCursor+1] = p.cursors[p.currentCursor]
	}
	p.currentCursor++
	return p.cursors[p.currentCursor]
}
func (p *Parser) DeleteCursor(err *error) {
	if (*err) != nil {
		p.currentCursor--
	}
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
		} else if !strings.EqualFold(string(c), string(otherChar)) {
			return false
		}
	}
	p.SetCursor(p.Cursor() + index + 1)
	return true
}

func (p *Parser) Regex(re *regexp.Regexp) (loc []int, err error) {
	cursor := p.SaveCursor()
	defer p.DeleteCursor(&err)
	if loc = re.FindReaderIndex(p); loc == nil {
		err = errors.New("expected regex match")
	} else {
		p.SetCursor(cursor + loc[1])
		loc[0] += cursor
		loc[1] += cursor
	}
	return
}

func (p *Parser) ReadRune() (r rune, size int, err error) {
	if r, err = p.CharacterAtIndex(p.Cursor()); err != nil {
		return -1, -1, err
	}
	size = utf8.RuneLen(r)
	p.IncrementCursor()
	return
}

func (p *Parser) MatchRange(loc []int) string {
	if loc[1] >= len(p.text) {
		return string(p.text[loc[0]:])
	}
	return string(p.text[loc[0]:loc[1]])
}

var literalRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_$]*`)

func (p *Parser) Literal() (literal string, acceptError error) {
	p.SaveCursor()
	defer p.DeleteCursor(&acceptError)
	if matchResult, err := p.Regex(literalRegex); err != nil {
		acceptError = errors.New("expected literal")
	} else {
		literal = p.MatchRange(matchResult)
	}
	return
}

func (p *Parser) AcceptContinuationToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	defer p.DeleteCursor(&acceptError)
	if !forest.JoinIgnoreError(p.AcceptNewLineToken()) {
		acceptError = errors.New("expected newline")
		return
	}
	forest.JoinIgnoreError(p.AcceptWhitespaceToken())
	if !p.Matches("1") {
		acceptError = errors.New("expected continuation literal '1'")
		return
	}
	forest.JoinIgnoreError(p.AcceptWhitespaceToken())
	return
}

var whitespaceRegex = regexp.MustCompile(`^[ \t]+`)

func (p *Parser) AcceptWhitespaceToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	defer p.DeleteCursor(&acceptError)
	if matchResult, err := p.Regex(whitespaceRegex); err != nil {
		if continuationForest, err := p.AcceptContinuationToken(); err == nil {
			forest.Join(continuationForest)
			return
		}
		acceptError = errors.New("expected whitespace")
	} else {
		forest.Add(WhitespaceToken{value: p.MatchRange(matchResult)})
		forest.JoinIgnoreError(p.AcceptContinuationToken())
	}
	return
}

var newLineRegex = regexp.MustCompile(`^[\r\n]+`)

func (p *Parser) AcceptNewLineToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	defer p.DeleteCursor(&acceptError)
	if matchResult, err := p.Regex(newLineRegex); err != nil {
		acceptError = errors.New("expected newline(s)")
	} else {
		forest.Add(NewLineToken{value: p.MatchRange(matchResult)})
	}
	return
}

var commentRegex = regexp.MustCompile(`^\![^\r\n]*`)

func (p *Parser) AcceptCommentToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	defer p.DeleteCursor(&acceptError)
	forest.JoinIgnoreError(p.AcceptWhitespaceToken())

	if matchResult, err := p.Regex(commentRegex); err != nil {
		acceptError = errors.New("expected comment")
	} else {
		forest.Add(CommentToken{value: p.MatchRange(matchResult)})
	}
	return
}

var stringRegex = regexp.MustCompile(`^[^'\r\n]*`)

func (p *Parser) AcceptStringToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	defer p.DeleteCursor(&acceptError)

	if !p.Matches("'") {
		acceptError = errors.New("expected left quote")
		return
	}

	if matchResult, err := p.Regex(stringRegex); err != nil {
		acceptError = errors.New("expected string value")
		return
	} else {
		token := StringToken{value: p.MatchRange(matchResult)}
		forest.Add(token)
	}

	if !p.Matches("'") {
		acceptError = errors.New("expected right quote")
		return
	}
	return
}

func (p *Parser) AcceptExpressionToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	defer p.DeleteCursor(&acceptError)

	if forest.JoinIgnoreError(p.AcceptStringToken()) {
		return
	}

	acceptError = errors.New("expected expression")
	return
}

func (p *Parser) AcceptPrimitiveCallToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	defer p.DeleteCursor(&acceptError)
	forest.JoinIgnoreError(p.AcceptWhitespaceToken())

	token := PrimitiveCallToken{}
	if p.Matches("print") {
		token.kind = "print"
	} else if p.Matches("type") {
		token.kind = "type"
	} else {
		acceptError = errors.New("expected primitive call statement")
		return
	}

	if _, err := p.AcceptWhitespaceToken(); err != nil {
		acceptError = errors.New("expected whitespace")
		return
	}

	if !p.Matches("*") {
		acceptError = errors.New("expected '*'")
		return
	}

	p.AcceptWhitespaceToken()

	if !p.Matches(",") {
		acceptError = errors.New("expected ','")
		return
	}

	p.AcceptWhitespaceToken()

	for token.values.JoinIgnoreError(p.AcceptExpressionToken()) {
		p.AcceptWhitespaceToken()
		if !p.Matches(",") {
			break
		}
		p.AcceptWhitespaceToken()
	}
	if token.values.Len() == 0 {
		acceptError = errors.New("expected 1 or more expressions")
		return
	}
	forest.Add(token)
	return
}

func (p *Parser) AcceptStatementTokens() (forest TokenForest, acceptError error) {
	for {
		if forest.JoinIgnoreError(p.AcceptWhitespaceToken()) {
			continue
		}
		if forest.JoinIgnoreError(p.AcceptNewLineToken()) {
			continue
		}
		if forest.JoinIgnoreError(p.AcceptCommentToken()) {
			continue
		}

		if forest.JoinIgnoreError(p.AcceptPrimitiveCallToken()) {
			continue
		}

		if forest.JoinIgnoreError(p.AcceptProgramToken()) {
			continue
		}

		return
	}

}

func (p *Parser) AcceptProgramToken() (forest TokenForest, acceptError error) {
	p.SaveCursor()
	defer p.DeleteCursor(&acceptError)
	if !p.Matches("program") {
		acceptError = errors.New("expected 'program'")
		return
	}
	if _, err := p.AcceptWhitespaceToken(); err != nil {
		acceptError = errors.New("expected whitespace")
		return
	}
	if name, err := p.Literal(); err != nil {
		acceptError = err
		return
	} else {
		token := ProgramToken{name: name, statements: CreateTokenForest()}
		token.statements.JoinIgnoreError(p.AcceptStatementTokens())
		forest.Add(token)
	}
	if p.Matches("endprogram") || p.Matches("end program") {
		return
	}
	acceptError = errors.New("did not find end program statement")
	return
}
