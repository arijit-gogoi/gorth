package word

type WordType string

type Token struct {
	Type    WordType
	Literal string
}

const (
	PUSH = "PUSH"
	ADD  = "ADD"
)
