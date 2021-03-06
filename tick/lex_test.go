package tick

import (
	"testing"
)

func TestLexer(t *testing.T) {

	type testCase struct {
		in     string
		tokens []token
	}

	test := func(tc testCase) {
		l := lex(tc.in)
		i := 0
		var tok token
		var ok bool
		for tok, ok = l.nextToken(); ok; tok, ok = l.nextToken() {
			if i >= len(tc.tokens) {
				t.Fatalf("unexpected extra token %v", tok)
			}
			if tok != tc.tokens[i] {
				t.Errorf("unexpected token: got %v exp %v i: %d in %s", tok, tc.tokens[i], i, tc.in)
			}
			i++
		}

		if i != len(tc.tokens) {
			t.Error("missing tokens", tc.tokens[i:])
		}
	}

	cases := []testCase{
		//Symbols + Operators
		{
			in: "!",
			tokens: []token{
				token{TokenNot, 0, "!"},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: "+",
			tokens: []token{
				token{TokenPlus, 0, "+"},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: "-",
			tokens: []token{
				token{TokenMinus, 0, "-"},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: "*",
			tokens: []token{
				token{TokenMult, 0, "*"},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: "/",
			tokens: []token{
				token{TokenDiv, 0, "/"},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: "=",
			tokens: []token{
				token{TokenAsgn, 0, "="},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: "==",
			tokens: []token{
				token{TokenEqual, 0, "=="},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "!=",
			tokens: []token{
				token{TokenNotEqual, 0, "!="},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: ">",
			tokens: []token{
				token{TokenGreater, 0, ">"},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: ">=",
			tokens: []token{
				token{TokenGreaterEqual, 0, ">="},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "<",
			tokens: []token{
				token{TokenLess, 0, "<"},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: "<=",
			tokens: []token{
				token{TokenLessEqual, 0, "<="},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "=~",
			tokens: []token{
				token{TokenRegexEqual, 0, "=~"},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "!~",
			tokens: []token{
				token{TokenRegexNotEqual, 0, "!~"},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "(",
			tokens: []token{
				token{TokenLParen, 0, "("},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: ")",
			tokens: []token{
				token{TokenRParen, 0, ")"},
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: ".",
			tokens: []token{
				token{TokenDot, 0, "."},
				token{TokenEOF, 1, ""},
			},
		},
		// Keywords
		{
			in: "AND",
			tokens: []token{
				token{TokenAnd, 0, "AND"},
				token{TokenEOF, 3, ""},
			},
		},
		{
			in: "OR",
			tokens: []token{
				token{TokenOr, 0, "OR"},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "TRUE",
			tokens: []token{
				token{TokenTrue, 0, "TRUE"},
				token{TokenEOF, 4, ""},
			},
		},
		{
			in: "FALSE",
			tokens: []token{
				token{TokenFalse, 0, "FALSE"},
				token{TokenEOF, 5, ""},
			},
		},
		{
			in: "var",
			tokens: []token{
				token{TokenVar, 0, "var"},
				token{TokenEOF, 3, ""},
			},
		},
		//Numbers
		{
			in: "42",
			tokens: []token{
				token{TokenNumber, 0, "42"},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "42.21",
			tokens: []token{
				token{TokenNumber, 0, "42.21"},
				token{TokenEOF, 5, ""},
			},
		},
		{
			in: ".421",
			tokens: []token{
				token{TokenDot, 0, "."},
				token{TokenNumber, 1, "421"},
				token{TokenEOF, 4, ""},
			},
		},
		{
			in: "0.421",
			tokens: []token{
				token{TokenNumber, 0, "0.421"},
				token{TokenEOF, 5, ""},
			},
		},
		//Durations
		{
			in: "42s",
			tokens: []token{
				token{TokenDuration, 0, "42s"},
				token{TokenEOF, 3, ""},
			},
		},
		{
			in: "42.21m",
			tokens: []token{
				token{TokenDuration, 0, "42.21m"},
				token{TokenEOF, 6, ""},
			},
		},
		{
			in: ".421h",
			tokens: []token{
				token{TokenDot, 0, "."},
				token{TokenDuration, 1, "421h"},
				token{TokenEOF, 5, ""},
			},
		},
		{
			in: "0.421s",
			tokens: []token{
				token{TokenDuration, 0, "0.421s"},
				token{TokenEOF, 6, ""},
			},
		},
		{
			in: "1u",
			tokens: []token{
				token{TokenDuration, 0, "1u"},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "1µ",
			tokens: []token{
				token{TokenDuration, 0, "1µ"},
				token{TokenEOF, 3, ""},
			},
		},
		{
			in: "1ms",
			tokens: []token{
				token{TokenDuration, 0, "1ms"},
				token{TokenEOF, 3, ""},
			},
		},
		{
			in: "1h",
			tokens: []token{
				token{TokenDuration, 0, "1h"},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "1d",
			tokens: []token{
				token{TokenDuration, 0, "1d"},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: "1w",
			tokens: []token{
				token{TokenDuration, 0, "1w"},
				token{TokenEOF, 2, ""},
			},
		},
		//Identifier
		{
			in: "variable",
			tokens: []token{
				token{TokenIdent, 0, "variable"},
				token{TokenEOF, 8, ""},
			},
		},
		{
			in: "myVar01",
			tokens: []token{
				token{TokenIdent, 0, "myVar01"},
				token{TokenEOF, 7, ""},
			},
		},
		// References
		{
			in: `""`,
			tokens: []token{
				token{TokenReference, 0, `""`},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: `"ref with spaces"`,
			tokens: []token{
				token{TokenReference, 0, `"ref with spaces"`},
				token{TokenEOF, 17, ""},
			},
		},
		{
			in: `"ref\""`,
			tokens: []token{
				token{TokenReference, 0, `"ref\""`},
				token{TokenEOF, 7, ""},
			},
		},
		//Strings
		{
			in: `''`,
			tokens: []token{
				token{TokenString, 0, `''`},
				token{TokenEOF, 2, ""},
			},
		},
		{
			in: `''''''`,
			tokens: []token{
				token{TokenString, 0, `''''''`},
				token{TokenEOF, 6, ""},
			},
		},
		{
			in: `'str'`,
			tokens: []token{
				token{TokenString, 0, `'str'`},
				token{TokenEOF, 5, ""},
			},
		},
		{
			in: `'str\''`,
			tokens: []token{
				token{TokenString, 0, `'str\''`},
				token{TokenEOF, 7, ""},
			},
		},
		{
			in: `'''s'tr'''`,
			tokens: []token{
				token{TokenString, 0, `'''s'tr'''`},
				token{TokenEOF, 10, ""},
			},
		},
		{
			in: `'''s\'tr'''`,
			tokens: []token{
				token{TokenString, 0, `'''s\'tr'''`},
				token{TokenEOF, 11, ""},
			},
		},
		{
			in: `'''str'''`,
			tokens: []token{
				token{TokenString, 0, `'''str'''`},
				token{TokenEOF, 9, ""},
			},
		},
		// Regex -- can only be lexed within context
		{
			in: `=~ //`,
			tokens: []token{
				token{TokenRegexEqual, 0, "=~"},
				token{TokenRegex, 3, "//"},
				token{TokenEOF, 5, ""},
			},
		},
		{
			in: `!~ //`,
			tokens: []token{
				token{TokenRegexNotEqual, 0, "!~"},
				token{TokenRegex, 3, "//"},
				token{TokenEOF, 5, ""},
			},
		},
		{
			in: `= //`,
			tokens: []token{
				token{TokenAsgn, 0, "="},
				token{TokenRegex, 2, "//"},
				token{TokenEOF, 4, ""},
			},
		},
		{
			in: `= /^((.*)[a-z]+\S{0,2})|cat\/\/$/`,
			tokens: []token{
				token{TokenAsgn, 0, "="},
				token{TokenRegex, 2, `/^((.*)[a-z]+\S{0,2})|cat\/\/$/`},
				token{TokenEOF, 33, ""},
			},
		},

		//Space
		{
			in: " ",
			tokens: []token{
				token{TokenEOF, 1, ""},
			},
		},
		{
			in: " \t\n",
			tokens: []token{
				token{TokenEOF, 3, ""},
			},
		},
		//Combinations
		{
			in: "var x = avg()",
			tokens: []token{
				token{TokenVar, 0, "var"},
				token{TokenIdent, 4, "x"},
				token{TokenAsgn, 6, "="},
				token{TokenIdent, 8, "avg"},
				token{TokenLParen, 11, "("},
				token{TokenRParen, 12, ")"},
				token{TokenEOF, 13, ""},
			},
		},
		{
			in: "var x = avg().parallel(4)x.groupby('cpu').window().period(10s)",
			tokens: []token{
				token{TokenVar, 0, "var"},
				token{TokenIdent, 4, "x"},
				token{TokenAsgn, 6, "="},
				token{TokenIdent, 8, "avg"},
				token{TokenLParen, 11, "("},
				token{TokenRParen, 12, ")"},
				token{TokenDot, 13, "."},
				token{TokenIdent, 14, "parallel"},
				token{TokenLParen, 22, "("},
				token{TokenNumber, 23, "4"},
				token{TokenRParen, 24, ")"},
				token{TokenIdent, 25, "x"},
				token{TokenDot, 26, "."},
				token{TokenIdent, 27, "groupby"},
				token{TokenLParen, 34, "("},
				token{TokenString, 35, "'cpu'"},
				token{TokenRParen, 40, ")"},
				token{TokenDot, 41, "."},
				token{TokenIdent, 42, "window"},
				token{TokenLParen, 48, "("},
				token{TokenRParen, 49, ")"},
				token{TokenDot, 50, "."},
				token{TokenIdent, 51, "period"},
				token{TokenLParen, 57, "("},
				token{TokenDuration, 58, "10s"},
				token{TokenRParen, 61, ")"},
				token{TokenEOF, 62, ""},
			},
		},
		//Comments
		{
			in: "var x = avg()\n// Comment all of this is ignored\nx.groupby('cpu')",
			tokens: []token{
				token{TokenVar, 0, "var"},
				token{TokenIdent, 4, "x"},
				token{TokenAsgn, 6, "="},
				token{TokenIdent, 8, "avg"},
				token{TokenLParen, 11, "("},
				token{TokenRParen, 12, ")"},
				token{TokenIdent, 48, "x"},
				token{TokenDot, 49, "."},
				token{TokenIdent, 50, "groupby"},
				token{TokenLParen, 57, "("},
				token{TokenString, 58, "'cpu'"},
				token{TokenRParen, 63, ")"},
				token{TokenEOF, 64, ""},
			},
		},
		{
			in: "var x = avg()\n// Comment all of this is ignored",
			tokens: []token{
				token{TokenVar, 0, "var"},
				token{TokenIdent, 4, "x"},
				token{TokenAsgn, 6, "="},
				token{TokenIdent, 8, "avg"},
				token{TokenLParen, 11, "("},
				token{TokenRParen, 12, ")"},
				token{TokenEOF, 47, ""},
			},
		},
	}

	for _, tc := range cases {
		test(tc)
	}
}
