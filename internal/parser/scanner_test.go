package parser

import (
	"strings"
	"testing"

	"github.com/0xch4z/selectr/internal/parser/ast"
	"github.com/0xch4z/selectr/internal/parser/token"
	"github.com/google/go-cmp/cmp"
)

func getNodesFromScanner(s *Scanner) []ast.Node {
	var nodes []ast.Node

	for {
		curr := s.Scan()
		nodes = append(nodes, curr)

		if curr.Tok == token.EOF {
			break
		}
	}
	return nodes
}

func TestScannerScan(t *testing.T) {
	for _, fixture := range []struct {
		content  string
		expected []ast.Node
	}{
		{
			content: "foo.bar",
			expected: []ast.Node{
				{
					Tok:      token.Ident,
					Lit:      "foo",
					StartPos: 0,
					EndPos:   3,
				},
				{
					Tok:      token.Dot,
					Lit:      ".",
					StartPos: 3,
					EndPos:   4,
				},
				{
					Tok:      token.Ident,
					Lit:      "bar",
					StartPos: 4,
					EndPos:   7,
				},
				{
					Tok:      token.EOF,
					Lit:      "\x00",
					StartPos: 7,
					EndPos:   8,
				},
			},
		},
		{
			content: "[3].object.nestedObject",
			expected: []ast.Node{
				{
					Tok:      token.LBracket,
					Lit:      "[",
					StartPos: 0,
					EndPos:   1,
				},
				{
					Tok:      token.Int,
					Lit:      "3",
					StartPos: 1,
					EndPos:   2,
				},
				{
					Tok:      token.RBracket,
					Lit:      "]",
					StartPos: 2,
					EndPos:   3,
				},
				{
					Tok:      token.Dot,
					Lit:      ".",
					StartPos: 3,
					EndPos:   4,
				},
				{
					Tok:      token.Ident,
					Lit:      "object",
					StartPos: 4,
					EndPos:   10,
				},
				{
					Tok:      token.Dot,
					Lit:      ".",
					StartPos: 10,
					EndPos:   11,
				},
				{
					Tok:      token.Ident,
					Lit:      "nestedObject",
					StartPos: 11,
					EndPos:   23,
				},
				{
					Tok:      token.EOF,
					Lit:      "\x00",
					StartPos: 23,
					EndPos:   24,
				},
			},
		},
		{
			content: "foo#~\\bar",
			expected: []ast.Node{
				{
					Tok:      token.Ident,
					Lit:      "foo",
					StartPos: 0,
					EndPos:   3,
				},
				{
					Tok:      token.Illegal,
					Lit:      "#",
					StartPos: 3,
					EndPos:   4,
				},
				{
					Tok:      token.Illegal,
					Lit:      "~",
					StartPos: 4,
					EndPos:   5,
				},
				{
					Tok:      token.Illegal,
					Lit:      "\\",
					StartPos: 5,
					EndPos:   6,
				},
				{
					Tok:      token.Ident,
					Lit:      "bar",
					StartPos: 6,
					EndPos:   9,
				},
				{
					Tok:      1,
					Lit:      "\x00",
					StartPos: 9,
					EndPos:   10,
				},
			},
		},

		{
			content: `
"test"`,
			expected: []ast.Node{
				{
					Tok:      token.WS,
					Lit:      "\n",
					StartPos: 0,
					EndPos:   1,
				},
				{
					Tok:      token.String,
					Lit:      "\"test\"",
					StartPos: 1,
					EndPos:   7,
				},
				{
					Tok:      token.EOF,
					Lit:      "\x00",
					StartPos: 7,
					EndPos:   8,
				},
			},
		},

		{
			content: "",
			expected: []ast.Node{
				{
					Tok:      token.EOF,
					Lit:      "\x00",
					StartPos: 0,
					EndPos:   1,
				},
			},
		},
	} {
		scanner := NewScanner(strings.NewReader(fixture.content))
		nodes := getNodesFromScanner(scanner)

		if !cmp.Equal(nodes, fixture.expected) {
			t.Errorf("`%s` was not scanned as expected:\n%s", fixture.content, cmp.Diff(fixture.expected, nodes))
		}
	}

}
