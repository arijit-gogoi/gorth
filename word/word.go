package word

type WordType string

type Token struct {
	Word    WordType
	Literal string
}

const (
	PUSH = "PUSH"
	ADD  = "ADD"
)
