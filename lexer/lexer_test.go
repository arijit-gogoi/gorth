package lexer

import (
	"reflect"
	"slices"
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

	l := New(input)
	for i, tt := range output {
		tok, _ := l.NextToken()
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
		expectedRecord  map[string][]word.Word
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
				{word.PUSH, "420", map[string][]word.Word{}},
				{word.DUP, "dup", map[string][]word.Word{}},
			},
		},
		{
			name:  "cr cr cr",
			input: `cr cr cr`,
			output: []expected{
				{word.CR, "cr", map[string][]word.Word{}},
				{word.CR, "cr", map[string][]word.Word{}},
				{word.CR, "cr", map[string][]word.Word{}},
			},
		},
		{
			name:  "LT and GT",
			input: `1 2 < -2 > -1 =`,
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
			name:  "and",
			input: `10 12 and`,
			output: []expected{
				{word.PUSH, "10", map[string][]word.Word{}},
				{word.PUSH, "12", map[string][]word.Word{}},
				{word.AND, "and", map[string][]word.Word{}},
			},
		},
		{
			name:  "test or with two numbers",
			input: `10 12 or`,
			output: []expected{
				{word.PUSH, "10", map[string][]word.Word{}},
				{word.PUSH, "12", map[string][]word.Word{}},
				{word.OR, "or", map[string][]word.Word{}},
			},
		},
		{
			name:  "invert: bitwise not",
			input: `1 invert`,
			output: []expected{
				{word.PUSH, "1", map[string][]word.Word{}},
				{word.INVERT, "invert", map[string][]word.Word{}},
			},
		},
		{
			name:  "udf: double",
			input: `: double dup + ;`,
			output: []expected{
				{
					expectedType: word.UDF,
					expectedLiteral: "double",
					expectedRecord: map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.ADD, "+"},
							{word.SEMICOLON, ";"},
						},
					},
				},
			},
		},
		{
			name:  "udf: square",
			input: `: double dup * ;`,
			output: []expected{
				{
					expectedType: word.UDF,
					expectedLiteral: "double",
					expectedRecord: map[string][]word.Word{
						"double": []word.Word{
							{word.DUP, "dup"},
							{word.MULTIPLY, "*"},
							{word.SEMICOLON, ";"},
						},
					},
				},
			},
		},
		{
			name:  "udf: half",
			input: `: half 2 swap / ;`,
			output: []expected{
				{
					expectedType: word.UDF,
					expectedLiteral: "half",
					expectedRecord: map[string][]word.Word{
						"half": []word.Word{
							{word.PUSH, "2"},
							{word.SWAP, "swap"},
							{word.DIVIDE, "/"},
							{word.SEMICOLON, ";"},
						},
					},
				},
			},
		},
	}
	for i, tc := range tests {
		l := New(tc.input)

		for _, o := range tc.output {
			tok, d := l.NextToken()
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
			t.Run(tc.name, func(t *testing.T) {
				if !reflect.DeepEqual(d, o.expectedRecord) {
					t.Fatalf("tests[%d] - record wrong. expected=%v, got=%v", i, o.expectedRecord, d)
				}
			})
		}
	}
}

func TestReadUDF(t *testing.T) {
	type test struct {
		name             string
		input            string
		expectedUDF      string
		expectedDefStack []word.Word
	}
	tests := []test{
		{
			name:             "udf infinite loop",
			input:            ": myudf",
			expectedUDF:      "myudf",
			expectedDefStack: []word.Word{},
		},
		{
			name:             "just a word, no defStack",
			input:            ": myword ;",
			expectedUDF:      "myword",
			expectedDefStack: []word.Word{
				{word.SEMICOLON, ";"},
			},
		},
		{
			name:             "udf double",
			input:            ": double dup + ;",
			expectedUDF:      "double",
			expectedDefStack: []word.Word{
				{word.DUP, "dup"},
				{word.ADD, "+"},
				{word.SEMICOLON, ";"},
			},
		},
		{
			name:             "udf square",
			input:            ": square dup * ;",
			expectedUDF:      "square",
			expectedDefStack: []word.Word{
				{word.DUP, "dup"},
				{word.MULTIPLY, "*"},
				{word.SEMICOLON, ";"},
			},
		},
	}
	for i, tc := range tests {
		l := New(tc.input)
		udf, defStack := l.readUDF()
		t.Run(tc.name, func(t *testing.T) {
			if udf != tc.expectedUDF {
				t.Fatalf("tests[%d] - udf string wrong. expected=%q, got=%q", i, tc.expectedUDF, udf)
			}
		})
		t.Run(tc.name, func(t *testing.T) {
			if !slices.Equal(defStack, tc.expectedDefStack) {
				t.Fatalf("tests[%d] - defStack wrong. expected=%v, got=%v", i, tc.expectedDefStack, defStack)
			}
		})
	}
}
