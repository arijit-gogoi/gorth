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
		{word.INT, "5"},
		{word.INT, "-10"},
		{word.ADD, "+"},
		{word.INT, "1"},
		{word.POP, "."},
		{word.INT, "1"},
		{word.SUBTRACT, "-"},
		{word.INT, "-1"},
		{word.MULTIPLY, "*"},
		{word.INT, "-6"},
		{word.DIVIDE, "/"},
		{word.DUP, "dup"},
		{word.DROP, "drop"},
		{word.INT, "2"},
		{word.SWAP, "swap"},
		{word.OVER, "over"},
		{word.SPIN, "spin"},
		{word.POP, "."},
		{word.POP, "."},
		{word.POP, "."},
		{word.INT, "97"},
		{word.EMIT, "emit"},
	}

	l := New(input, map[word.Word][]word.Word{})
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
		expectedDictionary map[word.Word][]word.Word
	}
	type test struct {
		name       string
		dictionary map[word.Word][]word.Word
		input      string
		output     []expected
	}
	tests := []test{
		{
			name:       "true false",
			dictionary: map[word.Word][]word.Word{},
			input:      `true false invert`,
			output: []expected{
				{
					expectedType:       word.TRUE,
					expectedLiteral:    "true",
					expectedDictionary: map[word.Word][]word.Word{},
				},
				{
					expectedType:       word.FALSE,
					expectedLiteral:    "false",
					expectedDictionary: map[word.Word][]word.Word{},
				},
				{
					expectedType:       word.INVERT,
					expectedLiteral:    "invert",
					expectedDictionary: map[word.Word][]word.Word{},
				},
			},
		},
		{
			name:       "mod",
			dictionary: map[word.Word][]word.Word{},
			input:      `5 5 mod`,
			output: []expected{
				{
					expectedType:       word.INT,
					expectedLiteral:    "5",
					expectedDictionary: map[word.Word][]word.Word{},
				},
				{
					expectedType:       word.INT,
					expectedLiteral:    "5",
					expectedDictionary: map[word.Word][]word.Word{},
				},
				{
					expectedType:       word.MOD,
					expectedLiteral:    "mod",
					expectedDictionary: map[word.Word][]word.Word{},
				},
			},
		},
		{
			name:       "%",
			input:      `5 5 %`,
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{
					expectedType:       word.INT,
					expectedLiteral:    "5",
					expectedDictionary: map[word.Word][]word.Word{},
				},
				{
					expectedType:       word.INT,
					expectedLiteral:    "5",
					expectedDictionary: map[word.Word][]word.Word{},
				},
				{
					expectedType:       word.MOD,
					expectedLiteral:    "%",
					expectedDictionary: map[word.Word][]word.Word{},
				},
			},
		},
		{
			name:       "dup a number",
			input:      `420 dup`,
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{
					expectedType:       word.INT,
					expectedLiteral:    "420",
					expectedDictionary: map[word.Word][]word.Word{},
				},
				{
					expectedType:       word.DUP,
					expectedLiteral:    "dup",
					expectedDictionary: map[word.Word][]word.Word{},
				},
			},
		},
		{
			name:       "cr cr cr",
			input:      `cr cr cr`,
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{word.CR, "cr", map[word.Word][]word.Word{}},
				{word.CR, "cr", map[word.Word][]word.Word{}},
				{word.CR, "cr", map[word.Word][]word.Word{}},
			},
		},
		{
			name:       "LT and GT",
			input:      `1 2 < -2 > -1 =`,
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{word.INT, "1", map[word.Word][]word.Word{}},
				{word.INT, "2", map[word.Word][]word.Word{}},
				{word.LT, "<", map[word.Word][]word.Word{}},
				{word.INT, "-2", map[word.Word][]word.Word{}},
				{word.GT, ">", map[word.Word][]word.Word{}},
				{word.INT, "-1", map[word.Word][]word.Word{}},
				{word.EQ, "=", map[word.Word][]word.Word{}},
			},
		},
		{
			name:       "and",
			input:      `10 12 and`,
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{word.INT, "10", map[word.Word][]word.Word{}},
				{word.INT, "12", map[word.Word][]word.Word{}},
				{word.AND, "and", map[word.Word][]word.Word{}},
			},
		},
		{
			name:       "test or with two numbers",
			input:      `10 12 or`,
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{word.INT, "10", map[word.Word][]word.Word{}},
				{word.INT, "12", map[word.Word][]word.Word{}},
				{word.OR, "or", map[word.Word][]word.Word{}},
			},
		},
		{
			name:       "invert: bitwise not",
			input:      `1 invert`,
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{word.INT, "1", map[word.Word][]word.Word{}},
				{word.INVERT, "invert", map[word.Word][]word.Word{}},
			},
		},
		{
			name:       "udf: double",
			input:      `: double dup + ;`,
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[word.Word][]word.Word{
						word.Word{word.UDF, "double"}: []word.Word{
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
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[word.Word][]word.Word{
						word.Word{word.UDF, "double"}: []word.Word{
							{word.DUP, "dup"}, {word.MULTIPLY, "*"},
						},
					},
				},
			},
		},
		{
			name:       "udf: half",
			input:      `: half 2 swap / ;`,
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[word.Word][]word.Word{
						word.Word{word.UDF, "half"}: []word.Word{
							{word.INT, "2"},
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
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{
					expectedType:       word.INT,
					expectedLiteral:    "1",
					expectedDictionary: map[word.Word][]word.Word{},
				},
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[word.Word][]word.Word{
						word.Word{word.UDF, "double"}: []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
				},
				{
					expectedType:    word.INT,
					expectedLiteral: "10",
					expectedDictionary: map[word.Word][]word.Word{
						word.Word{word.UDF, "double"}: []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
						},
					},
				},
				{
					expectedType:    word.UDF,
					expectedLiteral: "double",
					expectedDictionary: map[word.Word][]word.Word{
						word.Word{word.UDF, "double"}: []word.Word{
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
			dictionary: map[word.Word][]word.Word{},
			output: []expected{
				{
					expectedType:    word.DEFINE,
					expectedLiteral: ":",
					expectedDictionary: map[word.Word][]word.Word{
						word.Word{word.UDF, "buzz?"}: []word.Word{
							{word.INT, "5"},
							{word.MOD, "mod"},
							{word.INT, "0"},
							{word.EQ, "="},
							{word.IF, "if"},
							{word.INT, "2"},
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
				l.DefineWord()
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

func TestDefineWord(t *testing.T) {
	type test struct {
		name               string
		dictionary         map[word.Word][]word.Word
		input              string
		expectedDictionary map[word.Word][]word.Word
	}
	tests := []test{
		{
			name:       "udf infinite loop",
			input:      ": myudf",
			dictionary: map[word.Word][]word.Word{},
			expectedDictionary: map[word.Word][]word.Word{
				word.Word{word.UDF, "myudf"}: nil,
			},
		},
		{
			name:       "just a word, no defStack",
			input:      ": myword ;",
			dictionary: map[word.Word][]word.Word{},
			expectedDictionary: map[word.Word][]word.Word{
				word.Word{word.UDF, "myword"}: nil,
			},
		},
		{
			name:       "udf: double",
			input:      ": double dup + ;",
			dictionary: map[word.Word][]word.Word{},
			expectedDictionary: map[word.Word][]word.Word{
				word.Word{word.UDF, "double"}: []word.Word{
					{word.DUP, "dup"},
					{word.ADD, "+"},
				},
			},
		},
		{
			name:       "udf: square",
			input:      ": square dup * ;",
			dictionary: map[word.Word][]word.Word{},
			expectedDictionary: map[word.Word][]word.Word{
				word.Word{word.UDF, "square"}: []word.Word{
					{word.DUP, "dup"},
					{word.MULTIPLY, "*"},
				},
			},
		},
		{
			name:       "udf: the double UDF",
			input:      `: double dup + ; 10 double`,
			dictionary: map[word.Word][]word.Word{},
			expectedDictionary: map[word.Word][]word.Word{
				word.Word{word.UDF, "double"}: []word.Word{
					{word.DUP, "dup"},
					{word.ADD, "+"},
				},
			},
		},
		{
			name:       "udf: full sentence",
			input:      `: double dup + ; 10 double`,
			dictionary: map[word.Word][]word.Word{},
			expectedDictionary: map[word.Word][]word.Word{
				word.Word{word.UDF, "double"}: []word.Word{
					{word.DUP, "dup"},
					{word.ADD, "+"},
				},
			},
		},
	}
	for _, tc := range tests {
		l := New(tc.input, tc.dictionary)
		l.DefineWord()
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.expectedDictionary, l.Dictionary) {
				t.Fatalf("l.Dictionary wrong. expected=%v, got=%v", tc.expectedDictionary, l.Dictionary)
			}
		})
	}
}
