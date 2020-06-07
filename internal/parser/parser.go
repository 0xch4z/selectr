package parser

import (
	"io"

	"github.com/Charliekenney23/selectr/internal/parser/ast"
	"github.com/Charliekenney23/selectr/internal/parser/token"
)

// EOF signals the end of a file.
var EOF = rune(0)

// Parser represents a parser.
type Parser struct {
	s    *Scanner
	errs ErrorList
	pos  int
	buf  struct {
		node *ast.Node // last read node
		n    int       // buffer size
	}
}

// scan gets the next node from the underlying *Scanner or from the
// buffer.
func (p *Parser) scan() *ast.Node {
	p.pos++

	// if the buffer is not empty, return the buffered node and mark it
	// as empty.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.node
	}

	// read next node and save to buffer in case it needs to be unread
	// later.
	node := p.s.Scan()
	p.buf.node = &node
	return &node
}

// unscan retains the previously read token on the buffer to be
// processed later.
func (p *Parser) unscan() {
	p.pos--
	p.buf.n = 1
}

// expect asserts that the next scanned token is of the specified token
// kind.
func (p *Parser) expect(tok token.Token) *ast.Node {
	node := p.scan()

	if node.Tok != tok {
		p.errs.Push(errExpected(node.StartPos, tok))

		// return nil to signal an error
		return nil
	}
	return node
}

// parseAttributeExpr parses an attribute expression.
func (p *Parser) parseAttributeExpr() *ast.AttrExpr {
	var dot *ast.Node
	firstExpr := p.pos == 0

	if node := p.scan(); node.Tok == token.Dot {
		dot = node
	} else if firstExpr {
		// if this is the first expression, dot is optional.
		p.unscan()
	} else {
		p.errs.Push(errUnexpected(node.StartPos, node.Lit))

		// return nil to signal an error
		return nil
	}

	attr := p.expect(token.Ident)
	return &ast.AttrExpr{
		Dot:  dot,
		Attr: attr,
	}
}

// parseStringLit parses a string literal.
func (p *Parser) parseStringLit() *ast.StringLit {
	node := p.expect(token.String)
	return &ast.StringLit{
		Node: node,
	}
}

// parseIntLit parses an integer literal.
func (p *Parser) parseIntLit() *ast.IntLit {
	node := p.expect(token.Int)
	return &ast.IntLit{
		Node: node,
	}
}

// parseLitExpr parses a literal expression.
func (p *Parser) parseLitExpr() ast.LitExpr {
	node := p.scan()

	switch node.Tok {
	case token.String:
		p.unscan()
		return p.parseStringLit()

	case token.Int:
		p.unscan()
		return p.parseIntLit()
	}

	p.errs.Push(errUnexpected(node.StartPos, node.Lit))
	return nil
}

// parseIndexExpr parses an index expression.
func (p *Parser) parseIndexExpression() *ast.IndexExpr {
	lbrack := p.expect(token.LBracket)
	if lbrack == nil {
		return nil
	}

	litExpr := p.parseLitExpr()
	if litExpr == nil {
		return nil
	}

	rbrack := p.expect(token.RBracket)
	if rbrack == nil {
		return nil
	}

	return &ast.IndexExpr{
		LBracket: lbrack,
		Index:    litExpr,
		RBracket: rbrack,
	}
}

// Parse parses a selector.
func (p *Parser) Parse() (exprs []ast.Expr, err error) {
ParseLoop:
	for {
		node := p.scan()
		var expr ast.Expr

		switch node.Tok {
		case token.EOF:
			// the selector has been terminated, we can stop parsing.
			break ParseLoop

		case token.WS:
			// whitespace does not yield an expression; ignore it.
			continue

		case token.Dot, token.Ident:
			p.unscan()
			expr = p.parseAttributeExpr()

		case token.LBracket:
			p.unscan()
			expr = p.parseIndexExpression()

		default:
			// attribute and index expression are the only valid top level
			// expressions.
			return nil, errUnexpected(node.StartPos, node.Lit)
		}

		// expr is only nil when an error has occurred that is captured
		// on the parser or scanner.
		if expr != nil {
			exprs = append(exprs, expr)
		}

		// throw any underlying scanner errors, if there are any.
		if len(p.s.errs) != 0 {
			return nil, p.s.errs
		}

		if len(p.errs) != 0 {
			return nil, p.errs
		}
	}

	return
}

// New returns a new instance of Parser.
func New(r io.Reader) *Parser {
	return &Parser{
		s: NewScanner(r),
	}
}
