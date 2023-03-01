package parser

import (
	"fmt"

	"github.com/0xch4z/selectr/internal/parser/token"
)

// Error represents a parser error.
type Error struct {
	Pos int
	Msg string
}

// Error implements (error).Error
func (e Error) Error() string {
	return e.Msg
}

// ErrorList represents a list of parser errors.
type ErrorList []*Error

// ErrorList implements (error).Error
func (l ErrorList) Error() string {
	switch len(l) {
	case 0:
		return "no errors"
	case 1:
		return l[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", l[0], len(l))
}

// Add adds a new parser error to the collection.
func (p *ErrorList) Push(err *Error) {
	*p = append(*p, err)
}

// errExpected returns an error signaling that a token of the given kind
// was expected.
func errExpected(pos int, typ token.Token) *Error {
	return &Error{
		Pos: pos,
		Msg: "expected " + token.Tokens[typ],
	}
}

// errUnexpected returns an error signaling that the token read was not
// expected.
func errUnexpected(pos int, lit string) *Error {
	return &Error{
		Pos: pos,
		Msg: fmt.Sprintf("unexpected token '%s'", lit),
	}
}
