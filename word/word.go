package word

type WordType int

type Word struct {
	Type    WordType
	Literal string
}

const (
	// Stack
	PUSH = iota
	POP
	DUP
	DROP
	SWAP
	OVER
	SPIN
	EMIT
	CR

	// Math Operations
	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE

	// Boolean Operations
	EQ
	LT
	GT
	NOTEQ
	AND
	OR
	INVERT

	// UDF
	UDF
	COLON
	SEMICOLON

	// extra
	NEWLINE
	EOF
	ILLEGAL
)

var Table = map[string]WordType{
	"+":      ADD,
	"*":      MULTIPLY,
	"-":      SUBTRACT,
	"/":      DIVIDE,
	".":      POP,
	"dup":    DUP,
	"drop":   DROP,
	"swap":   SWAP,
	"over":   OVER,
	"spin":   SPIN,
	"emit":   EMIT,
	"cr":     CR,
	"=":      EQ,
	"<":      LT,
	">":      GT,
	"!=":     NOTEQ,
	"and":    AND,
	"or":     OR,
	"invert": INVERT,
	":":      COLON,
	";":	  SEMICOLON,
}

func GetWordType(s string, dictionary map[string][]Word) WordType {
	if wT, ok := Table[s]; ok {
		return wT
	} else if _, ok := dictionary[s]; ok {
		return UDF
	}
	return ILLEGAL
}

var Dictionary = make(map[Word][]Word)
