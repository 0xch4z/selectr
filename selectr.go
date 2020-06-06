package selectr

import (
	"fmt"
	"strings"

	"github.com/Charliekenney23/selectr/internal/parser"
	"github.com/Charliekenney23/selectr/internal/parser/ast"
)

// ResolveError represents an error that occured while resolving a
// value from a key-path.
type ResolveError struct {
	Err     error
	KeyPath string
}

// Error implements (error).Error
func (err ResolveError) Error() string {
	return err.KeyPath + " " + err.Err.Error()
}

// resolverError wraps an error to make a ResolveError with the given
// key-path.
func resolveError(keyPath string, err error) ResolveError {
	return ResolveError{
		KeyPath: keyPath,
		Err:     err,
	}
}

// TypeError signals an unexpected type.
type TypeError struct {
	ExpectedType string
	Value        interface{}
}

// Error implements (error).Error
func (err TypeError) Error() string {
	return fmt.Sprintf("TypeError: expected type %s but got %T", err.ExpectedType, err.Value)
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

// Resolve resolves a value from a map entry.
func (r *MapEntryResolver) Resolve(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case map[string]interface{}:
		return v[r.Key], nil
	default:
		return nil, TypeError{
			ExpectedType: "map[string]interface {}",
			Value:        v,
		}
	}
}

// Expression returns the corresponding *ast.Expr.
func (r *MapEntryResolver) Expression() ast.Expr {
	return r.Expr
}

// MapEntryResolver resolves a map.
var _ Resolver = (*MapEntryResolver)(nil)

// sliceIndex indexes a slice.
type SliceElementResolver struct {
	Index int
	Expr  *ast.IndexExpr
}

// Traverse indexes a slice.
func (r *SliceElementResolver) Resolve(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case []interface{}:
		if r.Index > len(v)-1 {
			return nil, fmt.Errorf("index out of range; index is %d but length is only %d", v, len(v))
		}
		return v[r.Index], nil
	default:
		return nil, TypeError{
			ExpectedType: "[]interface {}",
			Value:        v,
		}
	}
}

// Expression returns the corresponding *ast.Expr.
func (r *SliceElementResolver) Expression() ast.Expr {
	return r.Expr
}

// sliceIndex implements Traverser
var _ Resolver = (*SliceElementResolver)(nil)

// Parse parses a selector.
func Parse(s string) (*TraversalTreeNode, error) {
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
				Key:  e.Attribute.Lit,
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

	return head, nil
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
			return nil, resolveError("", err)
		}
		curr = curr.Child
	}
	return v, nil
}
