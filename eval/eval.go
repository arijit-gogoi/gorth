package eval

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Jorghy-Del/gorth/stack"
	"github.com/Jorghy-Del/gorth/word"
)

func Execute(words []word.Word) []int {
	var s stack.Stack
	for _, w := range words {
		switch w.Type {
		case word.EQ:
			if s.Pop() == s.Pop() {
				s.Push(-1)
			} else {
				s.Push(0)
			}
		case word.LT:
			v1 := s.Pop()
			v2 := s.Pop()
			if v2 < v1 {
				s.Push(-1)
			} else {
				s.Push(0)
			}
		case word.GT:
			v1 := s.Pop()
			v2 := s.Pop()
			if v2 > v1 {
				s.Push(-1)
			} else {
				s.Push(0)
			}
		case word.AND:
			s.Push(s.Pop() & s.Pop())
		case word.OR:
			s.Push(s.Pop() | s.Pop())
		case word.INVERT:
			s.Push(^s.Pop())
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
			fmt.Printf("%x of type %d is illegal.\n", w.Literal, w.Type)
		case word.PUSH:
			v, e := strconv.Atoi(w.Literal)
			if e != nil {
				log.Fatal(e)
			}
			s.Push(v)
		default:
			log.Fatalf("reached default %s (%T) has type %v\n", w.Literal, w.Literal, w.Type)
		}
	}
	return s.Stk
}
