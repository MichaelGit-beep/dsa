package main

import (
	"errors"
	"fmt"
)

type Node struct {
	Val   interface{}
	Left  *Node
	Right *Node
}

type Stack struct {
	Nodes []*Node
}

func (st *Stack) Push(node *Node) {
	st.Nodes = append(st.Nodes, node)
}

func (st *Stack) Pop() (*Node, error) {
	if len(st.Nodes) <= 0 {
		return &Node{}, errors.New("Can't pop from the empty stack")
	}
	n := st.Nodes[len(st.Nodes)-1]
	st.Nodes = st.Nodes[:len(st.Nodes)-1]
	if n.Right != nil {
		st.Push(n.Right)
	}
	if n.Left != nil {
		st.Push(n.Left)
	}
	return n, nil
}

func main() {
	a := Node{Val: "a"}
	b := Node{Val: "b"}
	c := Node{Val: "c"}
	d := Node{Val: "d"}
	e := Node{Val: "e"}
	f := Node{Val: "f"}
	a.Left = &b
	a.Right = &c
	b.Left = &d
	b.Right = &e
	c.Right = &f
	new_stack := Stack{}
	new_stack.Push(&a)
	for i := 0; i < 6; i++ {
		fmt.Println(new_stack.Pop())
	}
}
