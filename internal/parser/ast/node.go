package ast

import "github.com/0xch4z/selectr/internal/parser/token"

// Node represents an abstract syntax Node.
type Node struct {
	Tok      token.Token
	Lit      string
	StartPos int
	EndPos   int
}
