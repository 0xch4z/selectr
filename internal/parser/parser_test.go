package parser

import (
	"strings"
	"testing"

	"github.com/Charliekenney23/selectr/internal/parser/ast"
	"github.com/Charliekenney23/selectr/internal/parser/token"
	"github.com/google/go-cmp/cmp"
)

func TestParserParse(t *testing.T) {
	for _, fixture := range []struct {
		content  string
		expected []ast.Expr
		err      error
	}{
		{
			content: "foo[5]",
			expected: []ast.Expr{
				&ast.AttrExpr{
					Attribute: &ast.Node{
						Tok:      token.Ident,
						Lit:      "foo",
						StartPos: 0,
						EndPos:   3,
					},
				},
				&ast.IndexExpr{
					LBracket: &ast.Node{
						Tok:      token.LBracket,
						Lit:      "[",
						StartPos: 3,
						EndPos:   4,
					},
					Index: &ast.IntLit{
						Node: &ast.Node{
							Tok:      token.Int,
							Lit:      "5",
							StartPos: 4,
							EndPos:   5,
						},
					},
					RBracket: &ast.Node{
						Tok:      token.RBracket,
						Lit:      "]",
						StartPos: 5,
						EndPos:   6,
					},
				},
			},
		},

		{
			content: `["foo"]['bar'].bazz[9]`,
			expected: []ast.Expr{
				&ast.IndexExpr{
					LBracket: &ast.Node{
						Tok:      token.LBracket,
						Lit:      "[",
						StartPos: 0,
						EndPos:   1,
					},
					Index: &ast.StringLit{
						Node: &ast.Node{
							Tok:      token.String,
							Lit:      "\"foo\"",
							StartPos: 1,
							EndPos:   6,
						},
					},
					RBracket: &ast.Node{
						Tok:      token.RBracket,
						Lit:      "]",
						StartPos: 6,
						EndPos:   7,
					},
				},
				&ast.IndexExpr{
					LBracket: &ast.Node{
						Tok:      token.LBracket,
						Lit:      "[",
						StartPos: 7,
						EndPos:   8,
					},
					Index: &ast.StringLit{
						Node: &ast.Node{
							Tok:      token.String,
							Lit:      "'bar'",
							StartPos: 8,
							EndPos:   13,
						},
					},
					RBracket: &ast.Node{
						Tok:      token.RBracket,
						Lit:      "]",
						StartPos: 13,
						EndPos:   14,
					},
				},
				&ast.AttrExpr{
					Dot: &ast.Node{
						Tok:      token.Dot,
						Lit:      ".",
						StartPos: 14,
						EndPos:   15,
					},
					Attribute: &ast.Node{
						Tok:      token.Ident,
						Lit:      "bazz",
						StartPos: 15,
						EndPos:   19,
					},
				},
				&ast.IndexExpr{
					LBracket: &ast.Node{
						Tok:      token.LBracket,
						Lit:      "[",
						StartPos: 19,
						EndPos:   20,
					},
					Index: &ast.IntLit{
						Node: &ast.Node{
							Tok:      token.Int,
							Lit:      "9",
							StartPos: 20,
							EndPos:   21,
						},
					},
					RBracket: &ast.Node{
						Tok:      token.RBracket,
						Lit:      "]",
						StartPos: 21,
						EndPos:   22,
					},
				},
			},
		},

		{
			content: `.foo
.bar`,
			expected: []ast.Expr{
				&ast.AttrExpr{
					Dot: &ast.Node{
						Tok:      token.Dot,
						Lit:      ".",
						StartPos: 0,
						EndPos:   1,
					},
					Attribute: &ast.Node{
						Tok:      token.Ident,
						Lit:      "foo",
						StartPos: 1,
						EndPos:   4,
					},
				},
				&ast.AttrExpr{
					Dot: &ast.Node{
						Tok:      token.Dot,
						Lit:      ".",
						StartPos: 5,
						EndPos:   6,
					},
					Attribute: &ast.Node{
						Tok:      token.Ident,
						Lit:      "bar",
						StartPos: 6,
						EndPos:   9,
					},
				},
			},
		},

		{
			content:  "foo#5",
			err:      errUnexpected(3, "#"),
			expected: nil,
		},

		{
			content:  "",
			expected: nil,
		},
	} {
		parser := New(strings.NewReader(fixture.content))
		exprs, err := parser.Parse()

		if !cmp.Equal(fixture.err, err) {
			t.Errorf("parser error for `%s` was not as expected:\n%s", fixture.content, cmp.Diff(fixture.err, err))
		}

		if !cmp.Equal(exprs, fixture.expected) {
			t.Errorf("`%s` was not parsed as expected:\n%s", fixture.content, cmp.Diff(fixture.expected, exprs))
		}
	}
}
