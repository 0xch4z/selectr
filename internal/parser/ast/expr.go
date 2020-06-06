package ast

import "strconv"

// Expr represents an abstract expression.
type Expr interface {
	StartPos() int
	EndPos() int
	expr()
}

// LitExpr represents a literal expression.
type LitExpr interface {
	Expr
	Value() interface{}
}

// IndexExpr represents an index into a slice.
type IndexExpr struct {
	LBracket *Node
	Index    LitExpr
	RBracket *Node
}

func (e *IndexExpr) StartPos() int {
	return e.LBracket.StartPos
}

func (e *IndexExpr) EndPos() int {
	return e.RBracket.EndPos
}

func (IndexExpr) expr() {}

// IndexExpr implements Expr
var _ Expr = (*IndexExpr)(nil)

// IndexExpr represents an attribute selector.
type AttrExpr struct {
	Dot  *Node
	Attr *Node
}

func (e *AttrExpr) StartPos() int {
	return e.Dot.StartPos
}

func (e *AttrExpr) EndPos() int {
	return e.Attr.EndPos
}

func (AttrExpr) expr() {}

// AttributeExpr implements Expr
var _ Expr = (*AttrExpr)(nil)

type StringLit struct {
	Node *Node
}

func (l *StringLit) StartPos() int {
	return l.Node.StartPos
}

func (l *StringLit) EndPos() int {
	return l.Node.EndPos
}

func (l *StringLit) Value() interface{} {
	return l.Node.Lit[1 : len(l.Node.Lit)-1]
}

func (StringLit) expr() {}

// StringLit implements Expr
var _ Expr = (*StringLit)(nil)

// StringLit implements LitExpr
var _ LitExpr = (*StringLit)(nil)

type IntLit struct {
	Node *Node
}

func (l *IntLit) StartPos() int {
	return l.Node.StartPos
}

func (l *IntLit) EndPos() int {
	return l.Node.EndPos
}

func (l *IntLit) Value() interface{} {
	n, _ := strconv.Atoi(l.Node.Lit)
	return n
}

func (IntLit) expr() {}

// IntLit implements Expr
var _ Expr = (*IntLit)(nil)

// IntLit implements LitExpr
var _ LitExpr = (*IntLit)(nil)
