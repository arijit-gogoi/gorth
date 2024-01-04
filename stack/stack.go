package stack

import "log"

type Stack struct {
	Stk []int
}

func (s *Stack) Push(v int) {
	s.Stk = append(s.Stk, v)
}

func (s *Stack) Pop() (top int) {
	top = s.Stk[len(s.Stk)-1]
	s.Stk = s.Stk[:len(s.Stk)-1]
	return top
}

func (s *Stack) Len() int {
	return len(s.Stk)
}

func (s *Stack) Top() int {
	if len(s.Stk) < 0 {
		log.Fatal("len(s.Stk) is 0")
	}
	return s.Stk[len(s.Stk)-1]
}

func (s *Stack) Second() int {
	return s.Stk[len(s.Stk)-2]
}
