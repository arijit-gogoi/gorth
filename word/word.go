package word

type WordType string

type Word struct {
	Type    WordType
	Literal string
}

const (
	PUSH    = "PUSH"
	ADD     = "ADD"
	ILLEGAL = "ILLEGAL"
)
