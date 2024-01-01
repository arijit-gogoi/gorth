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

	// extra
	INT
	EOF
	ILLEGAL
)

var table = map[string]WordType{
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
}

func GetWordType(s string) WordType {
	if tok, ok := table[s]; ok {
		return tok
	}
	return -1
}
