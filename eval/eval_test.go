package eval

import (
	"slices"
	"testing"

	"github.com/Jorghy-Del/gorth/lexer"
	"github.com/Jorghy-Del/gorth/word"
)

func TestEvalTable(t *testing.T) {
	type expected struct {
		expectedType       word.WordType
		expectedLiteral    string
		expectedDictionary map[string][]word.Word
		expectedStk        []int
	}
	type test struct {
		name       string
		input      string
		dictionary map[string][]word.Word
		output     []expected
	}
	tests := []test{
		{
			name:       "EQ",
			input:      `8 8 = 4 =`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "8", map[string][]word.Word{}, []int{8}},
				{word.INT, "8", map[string][]word.Word{}, []int{8, 8}},
				{word.EQ, "=", map[string][]word.Word{}, []int{-1}},
				{word.INT, "4", map[string][]word.Word{}, []int{-1, 4}},
				{word.EQ, "=", map[string][]word.Word{}, []int{0}},
			},
		},
		{
			name:       "or",
			input:      `0 0 or -1 or 0 or -1 or 1 or`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "0", map[string][]word.Word{}, []int{0}},
				{word.INT, "0", map[string][]word.Word{}, []int{0, 0}},
				{word.OR, "or", map[string][]word.Word{}, []int{0}},
				{word.INT, "-1", map[string][]word.Word{}, []int{0, -1}},
				{word.OR, "or", map[string][]word.Word{}, []int{-1}},
				{word.INT, "0", map[string][]word.Word{}, []int{-1, 0}},
				{word.OR, "or", map[string][]word.Word{}, []int{-1}},
				{word.INT, "-1", map[string][]word.Word{}, []int{-1, -1}},
				{word.OR, "or", map[string][]word.Word{}, []int{-1}},
				{word.INT, "1", map[string][]word.Word{}, []int{-1, 1}},
				{word.OR, "or", map[string][]word.Word{}, []int{-1}},
			},
		},
		{
			name:       "and",
			input:      `-1 -1 and 0 and -1 and 0 and`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "-1", map[string][]word.Word{}, []int{-1}},
				{word.INT, "-1", map[string][]word.Word{}, []int{-1, -1}},
				{word.AND, "and", map[string][]word.Word{}, []int{-1}},
				{word.INT, "0", map[string][]word.Word{}, []int{-1, 0}},
				{word.AND, "and", map[string][]word.Word{}, []int{0}},
				{word.INT, "-1", map[string][]word.Word{}, []int{0, -1}},
				{word.AND, "and", map[string][]word.Word{}, []int{0}},
				{word.INT, "0", map[string][]word.Word{}, []int{0, 0}},
				{word.AND, "and", map[string][]word.Word{}, []int{0}},
			},
		},
		{
			name:       "invert true, then invert false",
			input:      `true invert invert`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.TRUE, "true", map[string][]word.Word{}, []int{-1}},
				{word.INVERT, "invert", map[string][]word.Word{}, []int{0}},
				{word.INVERT, "invert", map[string][]word.Word{}, []int{-1}},
			},
		},
		{
			name:       "modulo",
			input:      `8 3 mod 3 mod`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "8", map[string][]word.Word{}, []int{8}},
				{word.INT, "3", map[string][]word.Word{}, []int{8, 3}},
				{word.MOD, "mod", map[string][]word.Word{}, []int{2}},
				{word.INT, "3", map[string][]word.Word{}, []int{2, 3}},
				{word.MOD, "mod", map[string][]word.Word{}, []int{2}},
			},
		},
		{
			name:       "add one and minus one",
			input:      `1 -1 +`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "1", map[string][]word.Word{}, []int{1}},
				{word.INT, "-1", map[string][]word.Word{}, []int{1, -1}},
				{word.ADD, "+", map[string][]word.Word{}, []int{0}},
			},
		},
		{
			name:       "subtract two from one",
			input:      `2 1 -`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "2", map[string][]word.Word{}, []int{2}},
				{word.INT, "1", map[string][]word.Word{}, []int{2, 1}},
				{word.SUBTRACT, "-", map[string][]word.Word{}, []int{-1}},
			},
		},
		{
			name:       "dup a number",
			input:      `420 dup`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "420", map[string][]word.Word{}, []int{420}},
				{word.DUP, "dup", map[string][]word.Word{}, []int{420, 420}},
			},
		},
		{
			name:       "cr cr cr",
			input:      `cr cr cr`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.CR, "cr", map[string][]word.Word{}, []int{}},
				{word.CR, "cr", map[string][]word.Word{}, []int{}},
				{word.CR, "cr", map[string][]word.Word{}, []int{}},
			},
		},
		{
			name:       "1 2 3 cr cr cr",
			input:      `1 2 3 cr cr cr`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "1", map[string][]word.Word{}, []int{1}},
				{word.INT, "2", map[string][]word.Word{}, []int{1, 2}},
				{word.INT, "3", map[string][]word.Word{}, []int{1, 2, 3}},
				{word.CR, "cr", map[string][]word.Word{}, []int{1, 2, 3}},
				{word.CR, "cr", map[string][]word.Word{}, []int{1, 2, 3}},
				{word.CR, "cr", map[string][]word.Word{}, []int{1, 2, 3}},
			},
		},
		{
			name:       "Single character logical operations",
			input:      `1 2 < -2 > -1 =`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "1", map[string][]word.Word{}, []int{1}},
				{word.INT, "2", map[string][]word.Word{}, []int{1, 2}},
				{word.LT, "<", map[string][]word.Word{}, []int{-1}},
				{word.INT, "-2", map[string][]word.Word{}, []int{-1, -2}},
				{word.GT, ">", map[string][]word.Word{}, []int{-1}},
				{word.INT, "-1", map[string][]word.Word{}, []int{-1, -1}},
				{word.EQ, "=", map[string][]word.Word{}, []int{-1}},
			},
		},
		{
			name:       "and",
			input:      `10 12 and`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "10", map[string][]word.Word{}, []int{10}},
				{word.INT, "12", map[string][]word.Word{}, []int{10, 12}},
				{word.AND, "and", map[string][]word.Word{}, []int{8}},
			},
		},
		{
			name:       "test or with two numbers",
			input:      `10 12 or`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "10", map[string][]word.Word{}, []int{10}},
				{word.INT, "12", map[string][]word.Word{}, []int{10, 12}},
				{word.OR, "or", map[string][]word.Word{}, []int{14}},
			},
		},
		{
			name:       "invert: bitwise not",
			input:      `1 invert -1 * invert`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.INT, "1", map[string][]word.Word{}, []int{1}},
				{word.INVERT, "invert", map[string][]word.Word{}, []int{-2}},
				{word.INT, "-1", map[string][]word.Word{}, []int{-2, -1}},
				{word.MULTIPLY, "*", map[string][]word.Word{}, []int{2}},
				{word.INVERT, "invert", map[string][]word.Word{}, []int{-3}},
			},
		},
		{
			name:       "udf: full sentence",
			input:      `2 : double dup + ; 10 double double`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					word.INT, "2",
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
					word.INT, "10",
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
			name:       "udf: evaluate half",
			input:      `: half 2 swap / ; 100 half`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					word.DEFINE, ":",
					map[string][]word.Word{
						"half": []word.Word{
							{word.INT, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{},
				},
				{
					word.INT, "100",
					map[string][]word.Word{
						"half": []word.Word{
							{word.INT, "2"},
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
							{word.INT, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{50},
				},
			},
		},
		{
			name:       "udf: evaluate double then half",
			input:      `: double dup + ; : half 2 swap / ; 100 double half`,
			dictionary: map[string][]word.Word{},
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
							{word.INT, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{},
				},
				{
					word.INT, "100",
					map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
						"half": []word.Word{
							{word.INT, "2"},
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
							{word.INT, "2"},
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
							{word.INT, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
					[]int{100},
				},
			},
		},
		{
			name:       "test simple if",
			input:      `: isTruthy? if -1 else 0 then ; 10 isTruthy?`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					word.DEFINE, ":",
					map[string][]word.Word{
						"isTruthy?": []word.Word{
							{word.IF, "if"},
							{word.INT, "-1"},
							{word.ELSE, "else"},
							{word.INT, "0"},
							{word.THEN, "then"},
						},
					},
					[]int{},
				},
				{
					word.INT, "10",
					map[string][]word.Word{
						"isTruthy?": []word.Word{
							{word.IF, "if"},
							{word.INT, "-1"},
							{word.ELSE, "else"},
							{word.INT, "0"},
							{word.THEN, "then"},
						},
					},
					[]int{10},
				},
				{
					word.UDF, "isTruthy?",
					map[string][]word.Word{
						"isTruthy?": []word.Word{
							{word.IF, "if"},
							{word.INT, "-1"},
							{word.ELSE, "else"},
							{word.INT, "0"},
							{word.THEN, "then"},
						},
					},
					[]int{10, -1},
				},
			},
		},
		{
			name:       "test falsy if",
			input:      `: isFalsy? if -1 else 0 then ; 0 isFalsy?`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					word.DEFINE, ":",
					map[string][]word.Word{
						"isFalsy?": []word.Word{
							{word.IF, "if"},
							{word.INT, "-1"},
							{word.ELSE, "else"},
							{word.INT, "0"},
							{word.THEN, "then"},
						},
					},
					[]int{},
				},
				{
					word.INT, "0",
					map[string][]word.Word{
						"isFalsy?": []word.Word{
							{word.IF, "if"},
							{word.INT, "-1"},
							{word.ELSE, "else"},
							{word.INT, "0"},
							{word.THEN, "then"},
						},
					},
					[]int{0},
				},
				{
					word.UDF, "isFalsy?",
					map[string][]word.Word{
						"isFalsy?": []word.Word{
							{word.IF, "if"},
							{word.INT, "-1"},
							{word.ELSE, "else"},
							{word.INT, "0"},
							{word.THEN, "then"},
						},
					},
					[]int{0, -1},
				},
			},
		},
		{
			name:       "udf if: push 420",
			input:      `: buzz? 5 mod 0 = if 420 else 0 then ; 10 buzz?`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					word.DEFINE, ":",
					map[string][]word.Word{
						"buzz?": []word.Word{
							{word.INT, "5"},
							{word.MOD, "mod"},
							{word.INT, "0"},
							{word.EQ, "="},
							{word.IF, "if"},
							{word.INT, "420"},
							{word.ELSE, "else"},
							{word.INT, "0"},
							{word.THEN, "then"},
						},
					},
					[]int{},
				},
				{
					word.INT, "10",
					map[string][]word.Word{
						"buzz?": []word.Word{
							{word.INT, "5"},
							{word.MOD, "mod"},
							{word.INT, "0"},
							{word.EQ, "="},
							{word.IF, "if"},
							{word.INT, "420"},
							{word.ELSE, "else"},
							{word.INT, "0"},
							{word.THEN, "then"},
						},
					},
					[]int{10},
				},
				{
					word.UDF, "buzz?",
					map[string][]word.Word{
						"buzz?": []word.Word{
							{word.INT, "5"},
							{word.MOD, "mod"},
							{word.INT, "0"},
							{word.EQ, "="},
							{word.IF, "if"},
							{word.INT, "420"},
							{word.ELSE, "else"},
							{word.INT, "0"},
							{word.THEN, "then"},
						},
					},
					[]int{10, 420},
				},
			},
		},
	}
	for i, tc := range tests {
		l := lexer.New(tc.input, tc.dictionary)
		script := []word.Word{}


		for n, o := range tc.output {
			tok := l.NextToken()
			switch tok.Type {
			case word.DEFINE:
				l.DefineWord()
			case word.UDF:
				def := l.Dictionary[tok.Literal]
				isConditional := false
				for _, t := range def {
					if t.Type == word.IF {
						isConditional = true
						consequent := def[len(def)-4]
						alternate := def[len(def)-2]
						topType := script[len(script)-1].Type
						if int(topType) != word.FALSE { // if truthy
							script = append(script, consequent)
						} else if int(topType) == word.FALSE {
							script = append(script, alternate)
						}
					}
				}
				if !isConditional {
					script = append(script, def...)
				}
			default:
				script = append(script, tok)
			}
			got := Execute(script)

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
