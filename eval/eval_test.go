package eval

import (
	"slices"
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
	l := lexer.New(input, map[string][]word.Word{})
	words := []word.Word{}
	for i, tt := range output {
		tok := l.NextToken()
		words = append(words, tok)
		got := Execute(words)
		t.Run("single", func(t *testing.T) {
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
			}
		})
		t.Run("single", func(t *testing.T) {
			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
			}
		})
		t.Run("single", func(t *testing.T) {
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
		expectedType       word.WordType
		expectedLiteral    string
		expectedDictionary map[string][]word.Word
		expectedStk        []int
	}
	type test struct {
		name   string
		input  string
		output []expected
	}
	tests := []test{
		{
			name:  "modulo: ",
			input: `8 3 mod 3 mod`,
			output: []expected{
				{word.PUSH, "8", map[string][]word.Word{}, []int{8}},
				{word.PUSH, "3", map[string][]word.Word{}, []int{8, 3}},
				{word.MOD, "mod", map[string][]word.Word{}, []int{2}},
				{word.PUSH, "3", map[string][]word.Word{}, []int{2, 3}},
				{word.MOD, "mod", map[string][]word.Word{}, []int{2}},
			},
		},
		{
			name:  "add one and minus one",
			input: `1 -1 +`,
			output: []expected{
				{word.PUSH, "1", map[string][]word.Word{}, []int{1}},
				{word.PUSH, "-1", map[string][]word.Word{}, []int{1, -1}},
				{word.ADD, "+", map[string][]word.Word{}, []int{0}},
			},
		},
		{
			name:  "subtract two from one",
			input: `2 1 -`,
			output: []expected{
				{word.PUSH, "2", map[string][]word.Word{}, []int{2}},
				{word.PUSH, "1", map[string][]word.Word{}, []int{2, 1}},
				{word.SUBTRACT, "-", map[string][]word.Word{}, []int{-1}},
			},
		},
		{
			name:  "dup a number",
			input: `420 dup`,
			output: []expected{
				{word.PUSH, "420", map[string][]word.Word{}, []int{420}},
				{word.DUP, "dup", map[string][]word.Word{}, []int{420, 420}},
			},
		},
		{
			name:  "cr cr cr",
			input: `cr cr cr`,
			output: []expected{
				{word.CR, "cr", map[string][]word.Word{}, []int{}},
				{word.CR, "cr", map[string][]word.Word{}, []int{}},
				{word.CR, "cr", map[string][]word.Word{}, []int{}},
			},
		},
		{
			name:  "1 2 3 cr cr cr",
			input: `1 2 3 cr cr cr`,
			output: []expected{
				{word.PUSH, "1", map[string][]word.Word{}, []int{1}},
				{word.PUSH, "2", map[string][]word.Word{}, []int{1, 2}},
				{word.PUSH, "3", map[string][]word.Word{}, []int{1, 2, 3}},
				{word.CR, "cr", map[string][]word.Word{}, []int{1, 2, 3}},
				{word.CR, "cr", map[string][]word.Word{}, []int{1, 2, 3}},
				{word.CR, "cr", map[string][]word.Word{}, []int{1, 2, 3}},
			},
		},
		{
			name:  "Single character logical operations",
			input: `1 2 < -2 > -1 =`,
			output: []expected{
				{word.PUSH, "1", map[string][]word.Word{}, []int{1}},
				{word.PUSH, "2", map[string][]word.Word{}, []int{1, 2}},
				{word.LT, "<", map[string][]word.Word{}, []int{-1}},
				{word.PUSH, "-2", map[string][]word.Word{}, []int{-1, -2}},
				{word.GT, ">", map[string][]word.Word{}, []int{-1}},
				{word.PUSH, "-1", map[string][]word.Word{}, []int{-1, -1}},
				{word.EQ, "=", map[string][]word.Word{}, []int{-1}},
			},
		},
		{
			name:  "and",
			input: `10 12 and`,
			output: []expected{
				{word.PUSH, "10", map[string][]word.Word{}, []int{10}},
				{word.PUSH, "12", map[string][]word.Word{}, []int{10, 12}},
				{word.AND, "and", map[string][]word.Word{}, []int{8}},
			},
		},
		{
			name:  "test or with two numbers",
			input: `10 12 or`,
			output: []expected{
				{word.PUSH, "10", map[string][]word.Word{}, []int{10}},
				{word.PUSH, "12", map[string][]word.Word{}, []int{10, 12}},
				{word.OR, "or", map[string][]word.Word{}, []int{14}},
			},
		},
		{
			name:  "invert: bitwise not",
			input: `1 invert -1 * invert`,
			output: []expected{
				{word.PUSH, "1", map[string][]word.Word{}, []int{1}},
				{word.INVERT, "invert", map[string][]word.Word{}, []int{-2}},
				{word.PUSH, "-1", map[string][]word.Word{}, []int{-2, -1}},
				{word.MULTIPLY, "*", map[string][]word.Word{}, []int{2}},
				{word.INVERT, "invert", map[string][]word.Word{}, []int{-3}},
			},
		},
		{
			name:  "udf: full sentence",
			input: `2 : double dup + ; 10 double double`,
			output: []expected{
				{
					word.PUSH, "2",
					map[string][]word.Word{},
					[]int{2},
				},
				{
					word.DEFINE, ":",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
					[]int{2},
				},
				{
					word.PUSH, "10",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
					[]int{2, 10},
				},
				{
					word.UDF, "double",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
					[]int{2, 20},
				},
				{
					word.UDF, "double",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
					[]int{2, 40},
				},
			},
		},
		{
			name:  "udf: evaluate half",
			input: `: half 2 swap / ; 100 half`,
			output: []expected{
				{
					word.DEFINE, ":",
					map[string][]word.Word{
						"half": []word.Word{
							{word.PUSH, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{},
				},
				{
					word.PUSH, "100",
					map[string][]word.Word{
						"half": []word.Word{
							{word.PUSH, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{100},
				},
				{
					word.UDF, "half",
					map[string][]word.Word{
						"half": []word.Word{
							{word.PUSH, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{50},
				},
			},
		},
		{
			name:  "udf: evaluate double then half",
			input: `: double dup + ; : half 2 swap / ; 100 double half`,
			output: []expected{
				{
					word.DEFINE, ":",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
					[]int{},
				},
				{
					word.DEFINE, ":",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
						"half": []word.Word{
							{word.PUSH, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{},
				},
				{
					word.PUSH, "100",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
						"half": []word.Word{
							{word.PUSH, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{100},
				},
				{
					word.UDF, "double",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
						"half": []word.Word{
							{word.PUSH, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{200},
				},
				{
					word.UDF, "half",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
						"half": []word.Word{
							{word.PUSH, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{100},
				},
			},
		},
		{
			name:  "udf if: push 420",
			input: `: buzz? 5 mod 0 = if 420 then ;`,
			output: []expected{
				{
					word.DEFINE, ":",
					map[string][]word.Word{
						"buzz?": []word.Word{
							{word.PUSH, "5"},
							{word.MOD, "mod"},
							{word.PUSH, "0"},
							{word.EQ, "="},
							{word.IF, "if"},
							{word.PUSH, "420"},
							{word.THEN, "then"},
						},
					},
					[]int{},
				},
			},
		},
	}
	for i, tc := range tests {
		l := lexer.New(tc.input, map[string][]word.Word{})
		words := []word.Word{}
		for n, o := range tc.output {
			tok := l.NextToken()

			if tok.Type == word.DEFINE {
				isConditional, _, _ := l.ParseUDF()
				if isConditional {
				}
			} else if tok.Type == word.UDF {
				words = append(words, l.Dictionary[tok.Literal]...)
			} else {
				words = append(words, tok)
			}

			got := Execute(words)
			t.Run(tc.name, func(t *testing.T) {
				if tok.Type != o.expectedType {
					t.Fatalf("tests[%d, %d] - tokentype wrong. expected=%v, got=%v", i, n, o.expectedType, tok.Type)
				}
			})
			t.Run(tc.name, func(t *testing.T) {
				if tok.Literal != o.expectedLiteral {
					t.Fatalf("tests[%d, %d] - literal wrong. expected=%q, got=%q", i, n, o.expectedLiteral, tok.Literal)
				}
			})
			t.Run(tc.name, func(t *testing.T) {
				if !slices.Equal(o.expectedStk, got) {
					t.Fatalf("wrong evaluation. expected=%v, got=%v", o.expectedStk, got)
				}
			})
		}
	}
}
