package selectr

import (
	"fmt"
	"strings"

	"github.com/Charliekenney23/selectr/internal/parser"
	"github.com/Charliekenney23/selectr/internal/parser/ast"
)

// ResolveError represents an error that occured while resolving a
// value for a given selector.
type ResolveError struct {
	Msg  string
	Code string

	// Pos is the position in the corresponding key-path of the underlying
	// expression that the resolver was dervied from.
	Pos int
}

// Error implements (error).Error
func (err ResolveError) Error() string {
	if err.Code != "" {
		return err.Code + ": " + err.Msg
	}
	return err.Msg
}

// Resolver resolves a value from an object.
type Resolver interface {
	Resolve(interface{}) (interface{}, error)
	Expression() ast.Expr
}

// MapEntryResolver resolves a value from a map.
type MapEntryResolver struct {
	Key  string
	Expr ast.Expr
}

// Resolve resolves the value of an entry on the map.
func (r *MapEntryResolver) Resolve(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case map[string]interface{}:
		return v[r.Key], nil
	}
	return nil, ResolveError{
		Code: "TypeError",
		Msg:  fmt.Sprintf("cannot resolve attribute '%s' on type %T", r.Key, v),
		Pos:  r.Expr.StartPos(),
	}
}

// Expression returns the corresponding ast.Expr.
func (r *MapEntryResolver) Expression() ast.Expr {
	return r.Expr
}

// MapEntryResolver implements Resolver.
var _ Resolver = (*MapEntryResolver)(nil)

// SliceElementResolve resolves values from a slice.
type SliceElementResolver struct {
	Index int
	Expr  *ast.IndexExpr
}

// Resolve resolves the value of the element at the index on the slice.
func (r *SliceElementResolver) Resolve(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []interface{}:
		if r.Index > len(v)-1 {
			return nil, fmt.Errorf("index out of range; index is %d but length is only %d", r.Index, len(v))
		}
		return v[r.Index], nil
	}
	return nil, ResolveError{
		Code: "TypeError",
		Msg:  fmt.Sprintf("cannot resolve element '%d' on type %T", r.Index, v),
		Pos:  r.Expr.StartPos(),
	}
}

// Expression returns the corresponding ast.Expr.
func (r *SliceElementResolver) Expression() ast.Expr {
	return r.Expr
}

// SliceElementResolver implements Resolver.
var _ Resolver = (*SliceElementResolver)(nil)

// Parse parses a traversal tree from the selector string and returns
// a new Selector instance.
func Parse(s string) (*Selector, error) {
	exprs, err := parser.New(strings.NewReader(s)).Parse()
	if err != nil {
		return nil, err
	}

	var head, tail *TraversalTreeNode
	for _, expr := range exprs {
		var resolver Resolver

		switch e := expr.(type) {
		case *ast.AttrExpr:
			resolver = &MapEntryResolver{
				Key:  e.Attr.Lit,
				Expr: e,
			}

		case *ast.IndexExpr:
			switch indexExpr := e.Index.(type) {
			case *ast.IntLit:
				resolver = &SliceElementResolver{
					Index: indexExpr.Value().(int),
					Expr:  e,
				}

			case *ast.StringLit:
				resolver = &MapEntryResolver{
					Key:  indexExpr.Value().(string),
					Expr: e,
				}
			}
		}

		curr := &TraversalTreeNode{
			Resolver: resolver,
		}

		if head == nil {
			head = curr
		} else {
			tail.Child = curr
			curr.Parent = tail
		}

		tail = curr
	}

	return &Selector{
		tree: head,
	}, nil
}

// TraversalTreeNode represents a tree node responsible for traversing
// a given object.
type TraversalTreeNode struct {
	Resolver Resolver
	Parent   *TraversalTreeNode
	Child    *TraversalTreeNode
}

// Selector represents a value selection on an object.
type Selector struct {
	tree *TraversalTreeNode
}

// Resolve resolves the value at the specified key-path, if any, from the
// provided object. The root object must be an indexable type such as a
// Map `map[string]interface{}` or a Slice `[]interface{}`.
//
// All errors will be prefixed with the sub-key-path the error occured at.
//
// Example usage:
//
//     sel := Parse("test[0].foo")
//     sel.Resolve(map[string]interface{}{
//	       "test": []map[string]interface{}{
//		       {"foo": "bar"}
//         }
//     })
//
func (s *Selector) Resolve(v interface{}) (interface{}, error) {
	curr := s.tree
	for curr != nil {
		var err error
		if v, err = curr.Resolver.Resolve(v); err != nil {
			return nil, err
		}
		curr = curr.Child
	}
	return v, nil
}
