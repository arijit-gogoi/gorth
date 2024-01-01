package eval

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Jorghy-Del/gorth/stack"
	"github.com/Jorghy-Del/gorth/word"
)

func eval(words []word.Word) []int {
	var s stack.Stack
	for _, w := range words {
		switch w.Type {
		case word.EQ:
			v1 := s.Pop()
			v2 := s.Pop()
			if v2 == v1 {
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
		case word.ADD:
			v1 := s.Pop()
			v2 := s.Pop()
			s.Push(v1 + v2)
		case word.SUBTRACT:
			v1 := s.Pop()
			v2 := s.Pop()
			s.Push(v1 - v2)
		case word.MULTIPLY:
			v1 := s.Pop()
			v2 := s.Pop()
			s.Push(v1 * v2)
		case word.DIVIDE:
			v1 := s.Pop()
			v2 := s.Pop()
			s.Push(v1 / v2)
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
			n1 := s.Pop()
			n2 := s.Pop()
			n3 := s.Pop()
			s.Push(n2)
			s.Push(n3)
			s.Push(n1)
		case word.EMIT:
			n := s.Pop()
			fmt.Println(string(rune(n)))
		case word.CR:
			fmt.Println()
		case word.PUSH:
			v, e := strconv.Atoi(w.Literal)
			if e != nil {
				log.Fatal(e)
			}
			s.Push(v)
		default:
			fmt.Println("you reached default")
			log.Fatal("You reached DEFAULT")
		}
	}
	return s.Stk
}
