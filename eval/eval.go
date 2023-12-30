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
		case word.ADD:
			v1 := s.Pop()
			v2 := s.Pop()
			s.Push(v1 + v2)
		case word.POP:
			top := s.Pop()
			fmt.Println(top)
		case word.PUSH:
			v, e := strconv.Atoi(w.Literal)
			if e != nil {
				log.Fatal(e)
			}
			s.Push(v)
		}
	}
	return s.Stk
}
