package eval

import (
	"fmt"
	"log"
	"strconv"

	"errors"

	"github.com/Jorghy-Del/gorth/stack"
	"github.com/Jorghy-Del/gorth/word"
)

func Execute(tokens []word.Word) ([]int, error) {
	var s stack.Stack
	for _, t := range tokens {
		switch t.Type {
		case word.TRUE:
			s.Push(-1)
		case word.FALSE:
			s.Push(0)
		case word.AND:
			s.Push(s.Pop() & s.Pop())
		case word.OR:
			s.Push(s.Pop() | s.Pop())
		case word.INVERT:
			s.Push(^s.Pop())
		case word.EQ:
			if s.Pop() == s.Pop() {
				s.Push(int(word.TRUE))
			} else {
				s.Push(int(word.FALSE))
			}
		case word.LT:
			v1 := s.Pop()
			v2 := s.Pop()
			if v2 < v1 {
				s.Push(int(word.TRUE))
			} else {
				s.Push(int(word.FALSE))
			}
		case word.GT:
			v1 := s.Pop()
			v2 := s.Pop()
			if v2 > v1 {
				s.Push(int(word.TRUE))
			} else {
				s.Push(int(word.FALSE))
			}
		case word.ADD:
			s.Push(s.Pop() + s.Pop())
		case word.SUBTRACT:
			s.Push(s.Pop() - s.Pop())
		case word.MULTIPLY:
			s.Push(s.Pop() * s.Pop())
		case word.DIVIDE:
			s.Push(s.Pop() / s.Pop())
		case word.MOD:
			f := s.Pop()
			sec := s.Pop()
			s.Push(sec % f)
		case word.POP:
			top := s.Pop()
			fmt.Println(top)
		case word.DUP:
			top := s.Top()
			s.Push(top)
		case word.DROP:
			s.Pop()
		case word.SWAP:
			first := s.Pop()
			second := s.Pop()
			s.Push(first)
			s.Push(second)
		case word.OVER:
			sec := s.Second()
			s.Push(sec)
		case word.SPIN:
			n1, n2, n3 := s.Pop(), s.Pop(), s.Pop()
			s.Push(n2)
			s.Push(n3)
			s.Push(n1)
		case word.EMIT:
			n := s.Pop()
			fmt.Println(string(rune(n)))
		case word.CR:
			fmt.Println()
		case word.EOF:
			fmt.Println()
		case word.ILLEGAL:
			fmt.Printf("%x of type %d is illegal.\n", t.Literal, t.Type)
		case word.INT:
			v, e := strconv.Atoi(t.Literal)
			if e != nil {
				log.Fatal(e)
			}
			s.Push(v)
		default:
			log.Fatalf("reached default %s (%T) has type %v\n", t.Literal, t.Literal, t.Type)
			return s.Stk, errors.New("reached default.")
		}
	}
	return s.Stk, nil
}
