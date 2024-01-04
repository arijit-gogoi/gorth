package lexer

import (
	"github.com/Jorghy-Del/gorth/word"
	"strings"
)

type Lexer struct {
	input        string
	ch           byte
	position     int
	readPosition int
	Dictionary   map[string][]word.Word
}

func New(input string, dictionary map[string][]word.Word) *Lexer {
	l := &Lexer{
		input: input,
		Dictionary: dictionary,
	}
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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}


func (l *Lexer) NextToken() (tok word.Word, defStack []word.Word) {
	l.skipWhitespace()

	switch l.ch {
	case '-':
		p := l.peekChar()
		if isDigit(p) {
			tok.Type = word.PUSH
			l.readChar()
			tok.Literal = "-" + l.readNumber()
			return tok, defStack
		} else {
			w := string(l.ch)
			tok = newToken(word.GetWordType(w, l.Dictionary), w)
		}
	case ':', ';', '.', '+', '*', '/', '%', '<', '>', '=':
		w := string(l.ch)
		tok = newToken(word.GetWordType(w, l.Dictionary), w)
	case 0x00:
		tok = newToken(word.EOF, "0x00")
	default:
		if isLetter(l.ch) {
			w := l.readWord()
			tok = newToken(word.GetWordType(w, l.Dictionary), w)
		} else if isDigit(l.ch) {
			tok.Type = word.PUSH
			tok.Literal = l.readNumber()
			return tok, defStack
		} else {
			tok = newToken(word.ILLEGAL, string(l.ch))
		}
	}
	l.readChar()
	return tok, defStack
}

func newToken(wT word.WordType, literal string) word.Word {
	return word.Word{Type: wT, Literal: literal}
}

func (l *Lexer) ReadUDF() {
	l.readChar() // skip ':'
	l.skipWhitespace()
	udf := l.readWord()
	definitionStack := []word.Word{}
	for l.ch != 0x00 {
		tok, _ := l.NextToken()
		if tok.Type == word.SEMICOLON && tok.Literal == ";" {
			break
		}
		if strings.HasSuffix(tok.Literal, "?") {
			// tok.Type == word.IF
		}
		definitionStack = append(definitionStack, tok)
	}
	l.Dictionary[udf] = definitionStack
}

func (l *Lexer) readWord() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || ch == '_' || ch == '?'
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
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' || l.ch == 0xd || l.ch == 0xa {
		l.readChar()
	}
}
