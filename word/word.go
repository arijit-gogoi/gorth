package word

type WordType int

type Word struct {
	Type    WordType
	Literal string
}

const (
	// Boolean Operations
	TRUE  WordType = -1
	FALSE WordType = 0
	EQ    WordType = iota
	NOTEQ
	LT
	GT
	AND
	OR
	INVERT // 8

	// Stack
	INT
	POP
	DUP
	DROP
	SWAP
	OVER
	SPIN
	EMIT
	CR // 17

	// Math Operations
	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE
	MOD // 22

	// Conditionals
	IF
	ELSE
	THEN // 25

	// UDF
	UDF
	DEFINE
	SEMICOLON // 28

	// extra
	NEWLINE
	EOF
	ILLEGAL // 31
)

var Table = map[string]WordType{
	"+":      ADD,
	"*":      MULTIPLY,
	"-":      SUBTRACT,
	"/":      DIVIDE,
	".":      POP,
	"%":      MOD,
	"mod":    MOD,
	"dup":    DUP,
	"drop":   DROP,
	"swap":   SWAP,
	"over":   OVER,
	"spin":   SPIN,
	"emit":   EMIT,
	"cr":     CR,
	"true":   TRUE,
	"false":  FALSE,
	"=":      EQ,
	"<":      LT,
	">":      GT,
	"!=":     NOTEQ,
	"and":    AND,
	"or":     OR,
	"invert": INVERT,
	":":      DEFINE,
	";":      SEMICOLON,
	"if":     IF,
	"else":   ELSE,
	"then":   THEN,
}

func GetWordType(s string, dictionary map[Word][]Word) WordType {
	if wT, ok := Table[s]; ok {
		return wT
	} else if _, ok := dictionary[Word{Type: UDF, Literal: s}]; ok {
		return UDF
	}
	return ILLEGAL
}
