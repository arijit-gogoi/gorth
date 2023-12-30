package word

type WordType string

type Word struct {
	Type    WordType
	Literal string
}

const (
	PUSH    = "PUSH"
	POP     = "POP"
	ADD     = "ADD"
	INT     = "INT"
	ILLEGAL = "ILLEGAL"
)
