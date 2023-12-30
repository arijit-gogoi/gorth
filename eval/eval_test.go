package eval

import (
	"testing"

	"github.com/Jorghy-Del/gorth/word"
	"github.com/Jorghy-Del/gorth/lexer"
)

func TestEval(t *testing.T) {
	input := `5 10 + . 1`
	tests := []struct {
		expectedType    word.WordType
		expectedLiteral string
		expectedStk []int
	}{
		{word.PUSH, "5", []int{5}},
		{word.PUSH, "10", []int{5, 10}},
		{word.ADD, "+", []int{15}},
		{word.POP, ".", []int{}},
		{word.PUSH, "1", []int{1}},
	}
	l := lexer.New(input)
	words := []word.Word{}
	for i, tt := range tests {
		tok := l.NextToken()
		words = append(words, tok)
		got := eval(words)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Literal)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
		for j := range got {
			if got[j] != tt.expectedStk[j] {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		}
	}
}
