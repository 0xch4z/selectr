package selectr

import (
	"regexp"
	"testing"

	"github.com/Charliekenney23/selectr/internal/parser"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	cmpOmitMapEntryResolverExpr     = cmpopts.IgnoreFields(MapEntryResolver{}, "Expr")
	cmpOmitSliceElementResolverExpr = cmpopts.IgnoreFields(SliceElementResolver{}, "Expr")

	cmpOpts = []cmp.Option{cmpOmitMapEntryResolverExpr, cmpOmitSliceElementResolverExpr}
)

type parseTestFixture struct {
	selector string
	err      error
	errRegex *regexp.Regexp
	expected *TraversalTreeNode
}

func runParseTest(t *testing.T, fixture parseTestFixture) {
	t.Helper()

	sel, err := Parse(fixture.selector)

	if fixture.errRegex != nil {
		if err == nil {
			t.Error("expected error to match regex but none was thrown")
			return
		} else if !fixture.errRegex.Match([]byte(err.Error())) {
			t.Errorf("expected error to match pattern '%s' but got '%s'", fixture.errRegex, err.Error())
			return
		}
	} else if diff := cmp.Diff(fixture.err, err); diff != "" {
		t.Errorf("error for parsing `%s` was not as expected:\n%s", fixture.selector, diff)
		return
	}

	if sel == nil {
		if fixture.expected != nil {
			t.Errorf("`%s` unexpectedly resolved to a nil selector", fixture.selector)
		}
		return
	}

	if !cmp.Equal(sel.tree, fixture.expected, cmpOpts...) {
		t.Errorf("`%s` was not parsed as expected:\n%s", fixture.selector, cmp.Diff(fixture.expected, sel.tree, cmpOpts...))
	}
}

func treeFromResolverSequence(resolvers []Resolver) *TraversalTreeNode {
	var head, curr *TraversalTreeNode
	for _, resolver := range resolvers {
		node := &TraversalTreeNode{Resolver: resolver}
		if head == nil {
			head = node
		} else {
			curr.Child = node
			node.Parent = curr
		}
		curr = node
	}
	return head
}

func TestParse_empty(t *testing.T) {
	// an empty selector should resolve to an empty tree.
	runParseTest(t, parseTestFixture{
		selector: "",
		expected: nil,
	})

	// same with whitespace
	runParseTest(t, parseTestFixture{
		selector: " \t\n",
		expected: nil,
	})
}

func TestParse_error(t *testing.T) {
	// *parser.Error errors should be propegated.
	runParseTest(t, parseTestFixture{
		selector: "#illegal",
		err: &parser.Error{
			Pos: 0,
			Msg: "unexpected token '#'",
		},
	})

	runParseTest(t, parseTestFixture{
		selector: "\"astring\"",
		err: &parser.Error{
			Pos: 0,
			Msg: "unexpected token '\"astring\"'",
		},
	})

	runParseTest(t, parseTestFixture{
		selector: ".",
		errRegex: regexp.MustCompile("expected IDENT"),
	})

	runParseTest(t, parseTestFixture{
		selector: "[4[",
		errRegex: regexp.MustCompile("expected ]"),
	})

	runParseTest(t, parseTestFixture{
		selector: "['\n']",
		errRegex: regexp.MustCompile("unterminated string literal"),
	})
}

func TestParse(t *testing.T) {
	runParseTest(t, parseTestFixture{
		selector: "foo",
		expected: treeFromResolverSequence([]Resolver{
			&MapEntryResolver{
				Key: "foo",
			},
		}),
	})

	runParseTest(t, parseTestFixture{
		selector: "['bar']",
		expected: treeFromResolverSequence([]Resolver{
			&MapEntryResolver{
				Key: "bar",
			},
		}),
	})

	runParseTest(t, parseTestFixture{
		selector: "[40]",
		expected: treeFromResolverSequence([]Resolver{
			&SliceElementResolver{
				Index: 40,
			},
		}),
	})

	runParseTest(t, parseTestFixture{
		selector: ".foo['bar'][0].object.nestedObject[10]",
		expected: treeFromResolverSequence([]Resolver{
			&MapEntryResolver{
				Key: "foo",
			},
			&MapEntryResolver{
				Key: "bar",
			},
			&SliceElementResolver{
				Index: 0,
			},
			&MapEntryResolver{
				Key: "object",
			},
			&MapEntryResolver{
				Key: "nestedObject",
			},
			&SliceElementResolver{
				Index: 10,
			},
		}),
	})
}
