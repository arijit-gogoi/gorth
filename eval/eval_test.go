package eval

import (
	"testing"

	"github.com/Jorghy-Del/gorth/lexer"
	"github.com/Jorghy-Del/gorth/word"
)

func TestEval(t *testing.T) {
	input := `5 -10 + 1 . 1 - -1 * -6 / dup drop 2 swap over spin . . . 97 emit
1 2 < -2 > -1 =`
	type expected struct {
		expectedType    word.WordType
		expectedLiteral string
		expectedStk     []int
	}
	output := []expected{
		{word.PUSH, "5", []int{5}},
		{word.PUSH, "-10", []int{5, -10}},
		{word.ADD, "+", []int{-5}},
		{word.PUSH, "1", []int{-5, 1}},
		{word.POP, ".", []int{-5}},
		{word.PUSH, "1", []int{-5, 1}},
		{word.SUBTRACT, "-", []int{6}},
		{word.PUSH, "-1", []int{6, -1}},
		{word.MULTIPLY, "*", []int{-6}},
		{word.PUSH, "-6", []int{-6, -6}},
		{word.DIVIDE, "/", []int{1}},
		{word.DUP, "dup", []int{1, 1}},
		{word.DROP, "drop", []int{1}},
		{word.PUSH, "2", []int{1, 2}},
		{word.SWAP, "swap", []int{2, 1}},
		{word.OVER, "over", []int{2, 1, 2}},
		{word.SPIN, "spin", []int{1, 2, 2}},
		{word.POP, ".", []int{1, 2}},
		{word.POP, ".", []int{1}},
		{word.POP, ".", []int{}},
		{word.PUSH, "97", []int{97}},
		{word.EMIT, "emit", []int{}},
		{word.PUSH, "1", []int{1}},
		{word.PUSH, "2", []int{1, 2}},
		{word.LT, "<", []int{-1}},
		{word.PUSH, "-2", []int{-1, -2}},
		{word.GT, ">", []int{-1}},
		{word.PUSH, "-1", []int{-1, -1}},
		{word.EQ, "=", []int{-1}},
	}
	l := lexer.New(input)
	words := []word.Word{}
	for i, tt := range output {
		tok := l.NextToken()
		words = append(words, tok)
		got := Eval(words)
		t.Run("single", func(t *testing.T) {
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
			for j := range got {
				if got[j] != tt.expectedStk[j] {
					t.Fatalf("tests[%d] - wrong evaluation. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
				}
			}
		})
	}
}

func TestEvalTable(t *testing.T) {
	type expected struct {
		expectedType    word.WordType
		expectedLiteral string
		expectedStk     []int
	}
	type test struct {
		name   string
		input  string
		output []expected
	}
	tests := []test{
		{
			name:  "add one and minus one",
			input: `1 -1 +`,
			output: []expected{
				{word.PUSH, "1", []int{1}},
				{word.PUSH, "-1", []int{1, -1}},
				{word.ADD, "+", []int{0}},
			},
		},
		{
			name:  "subtract two from one",
			input: `2 1 -`,
			output: []expected{
				{word.PUSH, "2", []int{2}},
				{word.PUSH, "1", []int{2, 1}},
				{word.SUBTRACT, "-", []int{-1}},
			},
		},
		{
			name:  "dup a number",
			input: `420 dup`,
			output: []expected{
				{word.PUSH, "420", []int{420}},
				{word.DUP, "dup", []int{420, 420}},
			},
		},
		{
			name:  "cr cr cr",
			input: `cr cr cr`,
			output: []expected{
				{word.CR, "cr", []int{}},
				{word.CR, "cr", []int{}},
				{word.CR, "cr", []int{}},
			},
		},
		{
			name:  "1 2 3 cr cr cr",
			input: `1 2 3 cr cr cr`,
			output: []expected{
				{word.PUSH, "1", []int{1}},
				{word.PUSH, "2", []int{1, 2}},
				{word.PUSH, "3", []int{1, 2, 3}},
				{word.CR, "cr", []int{1, 2, 3}},
				{word.CR, "cr", []int{1, 2, 3}},
				{word.CR, "cr", []int{1, 2, 3}},
			},
		},
		{
			name:  "Single character logical operations",
			input: `1 2 < -2 > -1 =`,
			output: []expected{
				{word.PUSH, "1", []int{1}},
				{word.PUSH, "2", []int{1, 2}},
				{word.LT, "<", []int{-1}},
				{word.PUSH, "-2", []int{-1, -2}},
				{word.GT, ">", []int{-1}},
				{word.PUSH, "-1", []int{-1, -1}},
				{word.EQ, "=", []int{-1}},
			},
		},
		{
			name:  "and",
			input: `10 12 and`,
			output: []expected{
				{word.PUSH, "10", []int{10}},
				{word.PUSH, "12", []int{10, 12}},
				{word.AND, "and", []int{8}},
			},
		},
		{
			name:  "test or with two numbers",
			input: `10 12 or`,
			output: []expected{
				{word.PUSH, "10", []int{10}},
				{word.PUSH, "12", []int{10, 12}},
				{word.OR, "or", []int{14}},
			},
		},
		{
			name:  "invert: bitwise not",
			input: `1 invert -1 * invert`,
			output: []expected{
				{word.PUSH, "1", []int{1}},
				{word.INVERT, "invert", []int{-2}},
				{word.PUSH, "-1", []int{-2, -1}},
				{word.MULTIPLY, "*", []int{2}},
				{word.INVERT, "invert", []int{-3}},
			},
		},
	}

	for _, tc := range tests {
		l := lexer.New(tc.input)
		words := []word.Word{}

		for i, o := range tc.output {
			tok := l.NextToken()
			words = append(words, tok)
			got := Eval(words)
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
			t.Run(tc.name, func(t *testing.T) {
				for j := range got {
					if got[j] != o.expectedStk[j] {
						t.Fatalf("tests[%d] - wrong evaluation. expected=%q, got=%q", i, o.expectedLiteral, tok.Literal)
					}
				}
			})
		}
	}
}
