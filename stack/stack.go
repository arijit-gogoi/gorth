package stack

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
