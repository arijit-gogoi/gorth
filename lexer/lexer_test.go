package lexer

import (
	"github.com/Jorghy-Del/gorth/word"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `5 -10 + 1 . 1 - -1 * -6 / dup drop 2 swap over spin . . . 97 emit`

	output := []struct {
		expectedType    word.WordType
		expectedLiteral string
	}{
		{word.PUSH, "5"},
		{word.PUSH, "-10"},
		{word.ADD, "+"},
		{word.PUSH, "1"},
		{word.POP, "."},
		{word.PUSH, "1"},
		{word.SUBTRACT, "-"},
		{word.PUSH, "-1"},
		{word.MULTIPLY, "*"},
		{word.PUSH, "-6"},
		{word.DIVIDE, "/"},
		{word.DUP, "dup"},
		{word.DROP, "drop"},
		{word.PUSH, "2"},
		{word.SWAP, "swap"},
		{word.OVER, "over"},
		{word.SPIN, "spin"},
		{word.POP, "."},
		{word.POP, "."},
		{word.POP, "."},
		{word.PUSH, "97"},
		{word.EMIT, "emit"},
	}

	l := New(input)
	for i, tt := range output {
		tok := l.NextToken()
		t.Run("single", func(t *testing.T) {
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestNextTokenTable(t *testing.T) {
	type expected struct {
		expectedType    word.WordType
		expectedLiteral string
	}
	type test struct {
		name   string
		input  string
		output []expected
	}
	tests := []test{
		{
			name:  "dup a number",
			input: `420 dup`,
			output: []expected{
				{word.PUSH, "420"},
				{word.DUP, "dup"},
			},
		},
		{
			name:  "cr cr cr",
			input: `cr cr cr`,
			output: []expected{
				{word.CR, "cr"},
				{word.CR, "cr"},
				{word.CR, "cr"},
			},
		},
		{
			name:  "LT and GT",
			input: `1 2 < -2 > -1 =`,
			output: []expected{
				{word.PUSH, "1"},
				{word.PUSH, "2"},
				{word.LT, "<"},
				{word.PUSH, "-2"},
				{word.GT, ">"},
				{word.PUSH, "-1"},
				{word.EQ, "="},
			},
		},
		{
			name:  "and",
			input: `10 12 and`,
			output: []expected{
				{word.PUSH, "10"},
				{word.PUSH, "12"},
				{word.AND, "and"},
			},
		},
		{
			name:  "test or with two numbers",
			input: `10 12 or`,
			output: []expected{
				{word.PUSH, "10"},
				{word.PUSH, "12"},
				{word.OR, "or"},
			},
		},
		{
			name:  "invert: bitwise not",
			input: `1 invert`,
			output: []expected{
				{word.PUSH, "1"},
				{word.INVERT, "invert"},
			},
		},
	}
	for _, tc := range tests {
		l := New(tc.input)

		for i, o := range tc.output {
			tok := l.NextToken()
			t.Run(tc.name, func(t *testing.T) {
				if tok.Type != o.expectedType {
					t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, o.expectedType, tok.Type)
				}
			})
			t.Run(tc.name, func(t *testing.T) {
				if tok.Literal != o.expectedLiteral {
					t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, o.expectedLiteral, tok.Literal)
				}
			})
		}
	}
}
