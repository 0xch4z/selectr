package parser

import (
	"bufio"
	"bytes"
	"io"

	"github.com/Charliekenney23/selectr/internal/parser/ast"
	"github.com/Charliekenney23/selectr/internal/parser/token"
)

// Scanner represents a lexical scanner.
type Scanner struct {
	r    *bufio.Reader
	errs ErrorList
	pos  int
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the underlying *bufio.Reader and
// returns it's position.
//
// If the reader fails to read the next rune, an EOF is returned.
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return rune(0)
	}
	s.pos++
	return ch
}

// unread unreads the last read rune.
func (s *Scanner) unread() {
	// it's okay to swallow this error. (*bufio.Reader).UnreadRune only
	// throws an error if the last called method was not ReadRune.
	if err := s.r.UnreadRune(); err == nil {
		s.pos--
	}
}

// scanWhitespace consumes all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok token.Token, lit string) {
	var buf bytes.Buffer

	// Read every contiguous whitespace character into the buffer.
	// If a non-whitespace character is found, the loop will exit.
	for {
		if ch := s.read(); ch == EOF {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return token.WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok token.Token, lit string) {
	var buf bytes.Buffer

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == EOF {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return token.Ident, buf.String()
}

// scanInt consumes all contiguous digit runes.
func (s *Scanner) scanInt() (tok token.Token, lit string) {
	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == EOF {
			break
		} else if !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return token.Int, buf.String()
}

// scanString scans a string.
func (s *Scanner) scanString() (tok token.Token, lit string) {
	var buf bytes.Buffer

	// this is either `"` or `'`.
	quote := s.read()
	buf.WriteRune(quote)

ParseLoop:
	for {
		ch := s.read()
		switch ch {
		case '\n', EOF:
			s.errs.Push(&Error{Pos: s.pos, Msg: "unterminated string literal"})
			break ParseLoop

		case '\\':
			escapee := s.read()
			switch escapee {
			case quote, 'a', 'b', 'e', 'f', 'n', 'r', 't', 'v', '\\', '?':
				buf.WriteString("\\" + string(escapee))
			default:
				s.errs.Push(&Error{Pos: s.pos, Msg: "invalid escape sequence"})
				break ParseLoop
			}

		default:
			buf.WriteRune(ch)
		}

		if ch == quote {
			// string has been terminated
			break ParseLoop
		}
	}

	return token.String, buf.String()
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() ast.Node {
	startPos := s.pos
	ch := s.read()

	tok := token.Illegal
	lit := string(ch)

	// If we see whitespace, consume all contiguous whitespace.
	if isWhitespace(ch) {
		s.unread()
		tok, lit = s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		tok, lit = s.scanIdent()
	} else if isDigit(ch) {
		s.unread()
		tok, lit = s.scanInt()
	} else if isQoute(ch) {
		s.unread()
		tok, lit = s.scanString()
	}

	switch ch {
	case EOF:
		tok = token.EOF
	case '.':
		tok = token.Dot
	case '[':
		tok = token.LBracket
	case ']':
		tok = token.RBracket
	}

	return ast.Node{
		StartPos: startPos,
		EndPos:   startPos + len(lit),
		Tok:      tok,
		Lit:      lit,
	}
}
