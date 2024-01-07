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
	CR // 8

	// Math Operations
	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE
	MOD // 13

	// Boolean Operations
	EQ
	LT
	GT
	NOTEQ
	AND
	OR
	INVERT // 20

	// Conditionals
	IF
	ELSE
	THEN // 23

	// UDF
	UDF
	DEFINE
	SEMICOLON // 26

	// extra
	NEWLINE
	EOF
	ILLEGAL // 29
)

var Table = map[string]WordType{
	"+":      ADD,
	"*":      MULTIPLY,
	"-":      SUBTRACT,
	"/":      DIVIDE,
	".":      POP,
	"%":	  MOD,
	"mod":    MOD,
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
	":":      DEFINE,
	";":	  SEMICOLON,
	"if":	  IF,
	"else":   ELSE,
	"then":   THEN,
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
