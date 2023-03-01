package parser

import (
	"regexp"
	"strings"
	"testing"

	"github.com/0xch4z/selectr/internal/parser/ast"
	"github.com/0xch4z/selectr/internal/parser/token"
	"github.com/google/go-cmp/cmp"
)

type parserFixture struct {
	content  string
	expected []ast.Expr
	err      error
	errRegex *regexp.Regexp
}

func runParserTest(t *testing.T, fixture parserFixture) {
	t.Helper()

	parser := New(strings.NewReader(fixture.content))
	exprs, err := parser.Parse()

	if fixture.errRegex != nil {
		if !fixture.errRegex.Match([]byte(err.Error())) {
			t.Errorf("expected error to match pattern '%s' but got: '%s'", fixture.errRegex, err)
		}
	} else {
		if !cmp.Equal(fixture.err, err) {
			t.Errorf("parser error for `%s` was not as expected:\n%s", fixture.content, cmp.Diff(fixture.err, err))
		}
	}

	if !cmp.Equal(exprs, fixture.expected) {
		t.Errorf("`%s` was not parsed as expected:\n%s", fixture.content, cmp.Diff(fixture.expected, exprs))
	}
}

func TestParserParse_attributeExpressions(t *testing.T) {
	// simple attribute
	runParserTest(t, parserFixture{
		content: ".val",
		expected: []ast.Expr{
			&ast.AttrExpr{
				Dot: &ast.Node{
					Tok:      token.Dot,
					Lit:      ".",
					StartPos: 0,
					EndPos:   1,
				},
				Attr: &ast.Node{
					Tok:      token.Ident,
					Lit:      "val",
					StartPos: 1,
					EndPos:   4,
				},
			},
		},
	})

	// attribute with dot omitted
	runParserTest(t, parserFixture{
		content: "attr",
		expected: []ast.Expr{
			&ast.AttrExpr{
				Attr: &ast.Node{
					Tok:      token.Ident,
					Lit:      "attr",
					StartPos: 0,
					EndPos:   4,
				},
			},
		},
	})

	// throws error if IDENT token expectation fails following parsing
	// of a DOT.
	runParserTest(t, parserFixture{
		content: ".5",
		err:     ErrorList{errExpected(1, token.Ident)},
	})
}

func TestParserParse_indexExpressions(t *testing.T) {
	// simple array index
	runParserTest(t, parserFixture{
		content: "[5]",
		expected: []ast.Expr{
			&ast.IndexExpr{
				LBracket: &ast.Node{
					Tok:      token.LBracket,
					Lit:      "[",
					StartPos: 0,
					EndPos:   1,
				},
				Index: &ast.IntLit{
					Node: &ast.Node{
						Tok:      token.Int,
						Lit:      "5",
						StartPos: 1,
						EndPos:   2,
					},
				},
				RBracket: &ast.Node{
					Tok:      token.RBracket,
					Lit:      "]",
					StartPos: 2,
					EndPos:   3,
				},
			},
		},
	})

	// object index
	runParserTest(t, parserFixture{
		content: "[\"attr\"]",
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
						Lit:      "\"attr\"",
						StartPos: 1,
						EndPos:   7,
					},
				},
				RBracket: &ast.Node{
					Tok:      token.RBracket,
					Lit:      "]",
					StartPos: 7,
					EndPos:   8,
				},
			},
		},
	})

	// object index with single quotes
	runParserTest(t, parserFixture{
		content: "['test']",
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
						Lit:      "'test'",
						StartPos: 1,
						EndPos:   7,
					},
				},
				RBracket: &ast.Node{
					Tok:      token.RBracket,
					Lit:      "]",
					StartPos: 7,
					EndPos:   8,
				},
			},
		},
	})
}

func TestParserParseLitExpr_unexpectedToken(t *testing.T) {
	parser := New(strings.NewReader("{"))
	parser.parseLitExpr()

	expectedErrList := ErrorList{
		errUnexpected(0, "{"),
	}

	if diff := cmp.Diff(expectedErrList, parser.errs); diff != "" {
		t.Errorf("error from parsing illegal token in parseLitExpr was not as expected:\n%s", diff)
	}
}

func TestParserParse_empty(t *testing.T) {
	// should not return any ast.Expr from parsing an empty string
	runParserTest(t, parserFixture{
		content:  "",
		expected: nil,
	})

	// same for only whitespace
	runParserTest(t, parserFixture{
		content:  "\n\t\r",
		expected: nil,
	})
}

func TestParserParse_unexpectedTokenError(t *testing.T) {
	// should throw an unexpected token error when encountering
	// an illegal token.
	runParserTest(t, parserFixture{
		content:  "foo#5",
		err:      errUnexpected(3, "#"),
		expected: nil,
	})
}

func TestParserParse_expectedTokenError(t *testing.T) {
	// should throw an expected token error when a token assertion
	// fails.
	runParserTest(t, parserFixture{
		content: "[4",
		err:     ErrorList{errExpected(2, token.RBracket)},
	})
}

func TestParserParse_allowsWhitespace(t *testing.T) {
	runParserTest(t, parserFixture{
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
				Attr: &ast.Node{
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
				Attr: &ast.Node{
					Tok:      token.Ident,
					Lit:      "bar",
					StartPos: 6,
					EndPos:   9,
				},
			},
		},
	})
}

func TestParserParse_scannerError(t *testing.T) {
	// scanner errors should be propagated
	runParserTest(t, parserFixture{
		content:  `["\3"]`,
		errRegex: regexp.MustCompile("invalid escape sequence"),
	})
}

func TestParserParse(t *testing.T) {
	// complex selector
	runParserTest(t, parserFixture{
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
				Attr: &ast.Node{
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
	})
}
