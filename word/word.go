package word

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	PUSH = "PUSH"
	ADD  = "ADD"
)
