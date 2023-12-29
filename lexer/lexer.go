package lexer

import (
	"github.com/Jorghy-Del/gorth/word"
)

type Lexer struct {
	input        string
	ch           byte
	position     int
	readPosition int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() word.Word {
	var tok word.Word

	l.skipWhitespace()

	switch l.ch {
	case '+':
		tok = newToken(word.ADD, l.ch)
	default:
		if isDigit(l.ch) {
			tok.Type = word.PUSH
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(word.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func newToken(wordType word.WordType, ch byte) word.Word {
	return word.Word{Type: wordType, Literal: string(ch)}
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == 'n' {
		l.readChar()
	}
}
