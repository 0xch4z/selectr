package token

// Token represents a lexical token.
type Token int

// Token constants begin with Tok.
const (
	// special tokens
	Illegal Token = iota
	EOF           // End of File
	WS            // Whitespace

	Ident
	Dot
	LBracket
	RBracket

	// value tokens
	String
	Int
)

// Tokens maps Token constants to their string representations.
var Tokens = [...]string{
	Illegal:  "ILLEGAL",
	EOF:      "EOF",
	Ident:    "IDENT",
	Dot:      ".",
	LBracket: "[",
	RBracket: "]",
	String:   "STRING",
	Int:      "INT",
}
