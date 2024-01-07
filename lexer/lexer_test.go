package lexer

import (
	"reflect"
	"testing"

	"github.com/Jorghy-Del/gorth/word"
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

	l := New(input, map[string][]word.Word{})
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
		expectedType       word.WordType
		expectedLiteral    string
		expectedDictionary map[string][]word.Word
	}
	type test struct {
		name       string
		dictionary map[string][]word.Word
		input      string
		output     []expected
	}
	tests := []test{
		{
			name:       "mod",
			dictionary: map[string][]word.Word{},
			input:      `5 5 mod`,
			output: []expected{
				{
					expectedType:       word.PUSH,
					expectedLiteral:    "5",
					expectedDictionary: map[string][]word.Word{},
				},
				{
					expectedType:       word.PUSH,
					expectedLiteral:    "5",
					expectedDictionary: map[string][]word.Word{},
				},
				{
					expectedType:       word.MOD,
					expectedLiteral:    "mod",
					expectedDictionary: map[string][]word.Word{},
				},
			},
		},
		{
			name:       "%",
			input:      `5 5 %`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					expectedType:       word.PUSH,
					expectedLiteral:    "5",
					expectedDictionary: map[string][]word.Word{},
				},
				{
					expectedType:       word.PUSH,
					expectedLiteral:    "5",
					expectedDictionary: map[string][]word.Word{},
				},
				{
					expectedType:       word.MOD,
					expectedLiteral:    "%",
					expectedDictionary: map[string][]word.Word{},
				},
			},
		},
		{
			name:       "dup a number",
			input:      `420 dup`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					expectedType:       word.PUSH,
					expectedLiteral:    "420",
					expectedDictionary: map[string][]word.Word{},
				},
				{
					expectedType:       word.DUP,
					expectedLiteral:    "dup",
					expectedDictionary: map[string][]word.Word{},
				},
			},
		},
		{
			name:       "cr cr cr",
			input:      `cr cr cr`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.CR, "cr", map[string][]word.Word{}},
				{word.CR, "cr", map[string][]word.Word{}},
				{word.CR, "cr", map[string][]word.Word{}},
			},
		},
		{
			name:       "LT and GT",
			input:      `1 2 < -2 > -1 =`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.PUSH, "1", map[string][]word.Word{}},
				{word.PUSH, "2", map[string][]word.Word{}},
				{word.LT, "<", map[string][]word.Word{}},
				{word.PUSH, "-2", map[string][]word.Word{}},
				{word.GT, ">", map[string][]word.Word{}},
				{word.PUSH, "-1", map[string][]word.Word{}},
				{word.EQ, "=", map[string][]word.Word{}},
			},
		},
		{
			name:       "and",
			input:      `10 12 and`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.PUSH, "10", map[string][]word.Word{}},
				{word.PUSH, "12", map[string][]word.Word{}},
				{word.AND, "and", map[string][]word.Word{}},
			},
		},
		{
			name:       "test or with two numbers",
			input:      `10 12 or`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.PUSH, "10", map[string][]word.Word{}},
				{word.PUSH, "12", map[string][]word.Word{}},
				{word.OR, "or", map[string][]word.Word{}},
			},
		},
		{
			name:       "invert: bitwise not",
			input:      `1 invert`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{word.PUSH, "1", map[string][]word.Word{}},
				{word.INVERT, "invert", map[string][]word.Word{}},
			},
		},
		{
			name:       "udf: double",
			input:      `: double dup + ;`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
				},
			},
		},
		{
			name:       "udf: square",
			input:      `: double dup * ;`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"}, {word.MULTIPLY, "*"},
						},
					},
				},
			},
		},
		{
			name:       "udf: half",
			input:      `: half 2 swap / ;`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[string][]word.Word{
						"half": []word.Word{
							{word.PUSH, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
						},
					},
				},
			},
		},
		{
			name:       "udf: simple full sentence",
			input:      `1 : double dup + ; 10 double`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					expectedType:       word.PUSH,
					expectedLiteral:    "1",
					expectedDictionary: map[string][]word.Word{},
				},
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
				},
				{
					expectedType:    word.PUSH,
					expectedLiteral: "10",
					expectedDictionary: map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
				},
				{
					expectedType:    word.UDF,
					expectedLiteral: "double",
					expectedDictionary: map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
				},
			},
		},
		{
			name:       "udf if: push 2",
			input:      `: buzz? 5 mod 0 = if 2 then ;`,
			dictionary: map[string][]word.Word{},
			output: []expected{
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[string][]word.Word{
						"buzz?": []word.Word{
							{word.PUSH, "5"},
							{word.MOD, "mod"},
							{word.PUSH, "0"},
							{word.EQ, "="},
							{word.IF, "if"},
							{word.PUSH, "2"},
							{word.THEN, "then"},
						},
					},
				},
			},
		},
	}
	for i, tc := range tests {
		l := New(tc.input, tc.dictionary)
		for _, o := range tc.output {
			tok := l.NextToken()
			if tok.Type == word.DEFINE {
				l.ParseUDF()
			}
			t.Run(tc.name, func(t *testing.T) {
				if !reflect.DeepEqual(l.Dictionary, o.expectedDictionary) {
					t.Fatalf("l.Dictionary wrong. expected=%v, got=%v", o.expectedDictionary, l.Dictionary)
				}
			})
			t.Run(tc.name, func(t *testing.T) {
				if tok.Type != o.expectedType {
					t.Fatalf("tests[%d] - tokentype wrong. expected=%d, got=%d", i, o.expectedType, tok.Type)
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

func TestParseUDF(t *testing.T) {
	type test struct {
		name               string
		dictionary         map[string][]word.Word
		input              string
		expectedDictionary map[string][]word.Word
	}
	tests := []test{
		{
			name:       "udf infinite loop",
			input:      ": myudf",
			dictionary: map[string][]word.Word{},
			expectedDictionary: map[string][]word.Word{
				"myudf": nil,
			},
		},
		{
			name:       "just a word, no defStack",
			input:      ": myword ;",
			dictionary: map[string][]word.Word{},
			expectedDictionary: map[string][]word.Word{
				"myword": nil,
			},
		},
		{
			name:       "udf: double",
			input:      ": double dup + ;",
			dictionary: map[string][]word.Word{},
			expectedDictionary: map[string][]word.Word{
				"double": []word.Word{
					{word.DUP, "dup"},
					{word.ADD, "+"},
				},
			},
		},
		{
			name:       "udf: square",
			input:      ": square dup * ;",
			dictionary: map[string][]word.Word{},
			expectedDictionary: map[string][]word.Word{
				"square": []word.Word{
					{word.DUP, "dup"},
					{word.MULTIPLY, "*"},
				},
			},
		},
		{
			name:       "udf: the double UDF",
			input:      `: double dup + ; 10 double`,
			dictionary: map[string][]word.Word{},
			expectedDictionary: map[string][]word.Word{
				"double": []word.Word{
					{word.DUP, "dup"},
					{word.ADD, "+"},
				},
			},
		},
		{
			name:       "udf: full sentence",
			input:      `: double dup + ; 10 double`,
			dictionary: map[string][]word.Word{},
			expectedDictionary: map[string][]word.Word{
				"double": []word.Word{
					{word.DUP, "dup"},
					{word.ADD, "+"},
				},
			},
		},
	}
	for _, tc := range tests {
		l := New(tc.input, tc.dictionary)
		l.ParseUDF()
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.expectedDictionary, l.Dictionary) {
				t.Fatalf("l.Dictionary wrong. expected=%v, got=%v", tc.expectedDictionary, l.Dictionary)
			}
		})
	}
}
