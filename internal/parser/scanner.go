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
		// decrement the current position if the last read rune was
		// successfully unread
		s.pos--
	}
}

// scanWhitespace consumes all contiguous whitespace runes.
func (s *Scanner) scanWhitespace() (tok token.Token, lit string) {
	var buf bytes.Buffer

	// read every contiguous whitespace character into the buffer.
	// if a non-whitespace character or EOF occurs, the loop will exit.
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

	// read every contiguous ident character into the buffer.
	// if a non-alphanumeric character other than '_' or EOF occurs,
	// the loop will exit.
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

// scanInt consumes all contiguous integer runes.
func (s *Scanner) scanInt() (tok token.Token, lit string) {
	var buf bytes.Buffer

	// read every contiguous numeric character into the buffer.
	// if a non-numeric character or EOF occurs, the loop will exit.
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

	// we can be sure that quote is either a `'` or a `"` rune as this
	// method is only called in Scan if one of these quotes are detected.
	// the quote must be saved so we can check that it occurs again,
	// terminating the string.
	quote := s.read()
	buf.WriteRune(quote)

ParseLoop:
	for {
		ch := s.read()

		switch ch {
		case '\n', EOF:
			// if a new line or EOF occurs in the middle of a string literal
			// the string is invalid as it has not been terminated with an
			// end-quote.
			s.errs.Push(&Error{Pos: s.pos, Msg: "unterminated string literal"})
			break ParseLoop

		case '\\':
			escapee := s.read()
			switch escapee {
			case quote, 'a', 'b', 'e', 'f', 'n', 'r', 't', 'v', '\\', '?':
				buf.WriteString("\\" + string(escapee))
			default:
				// the specified string escape is not recognized, so it's invalid.
				// the position included with the error should be one back to the
				// start of the escape sequence with the backslash character.
				s.errs.Push(&Error{Pos: s.pos - 1, Msg: "invalid escape sequence"})
				break ParseLoop
			}

		default:
			// read subsequent string character into the buffer.
			buf.WriteRune(ch)
		}

		if ch == quote {
			// if the character matches the quote that started the string, then the
			// string is terminated and we can stop parsing the string.
			break ParseLoop
		}
	}

	return token.String, buf.String()
}

// Scan reads the next token.
func (s *Scanner) Scan() ast.Node {
	startPos := s.pos
	ch := s.read()

	tok := token.Illegal
	lit := string(ch)

	// parse multi character token types if detected.
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

	// set the token type of any single character token types.
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
