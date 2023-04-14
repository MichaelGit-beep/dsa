package main

import (
	"fmt"
)

type Node struct {
	Val   interface{}
	Left  *Node
	Right *Node
}

type Queue struct {
	Nodes []*Node
}

func (q *Queue) Push(n *Node) {
	q.Nodes = append([]*Node{n}, q.Nodes...)

}

func (q *Queue) Pull() *Node {
	if len(q.Nodes) <= 0 {
		return &Node{}
	}

	n := q.Nodes[len(q.Nodes)-1]
	q.Nodes = q.Nodes[:len(q.Nodes)-1]
	return n
}

func breadthFirstValues(r *Node) []interface{} {
	results := make([]interface{}, 0)
	if r == nil || r.Val == nil {
		return results
	}

	nodeQueue := Queue{}
	nodeQueue.Push(r)

	current := nodeQueue.Pull()
	for current.Val != nil {
		results = append(results, current.Val)
		if current.Left != nil {
			nodeQueue.Push(current.Left)
		}
		if current.Right != nil {
			nodeQueue.Push(current.Right)
		}
		current = nodeQueue.Pull()

	}

	return results
}

func treeSum(root *Node) int {
	sum := 0
	queue := Queue{Nodes: []*Node{root}}
	current := queue.Pull()
	for current.Val != nil {
		sum += current.Val.(int)
		if current.Left != nil {
			queue.Push(current.Left)
		}
		if current.Right != nil {
			queue.Push(current.Right)
		}
		current = queue.Pull()
	}
	return sum
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
	//       a
	//     /   \
	//   b      c
	//  / \      \
	// d   e	  f
	fmt.Println(breadthFirstValues(&a))
	n1 := Node{Val: 3}
	n2 := Node{Val: 11}
	n3 := Node{Val: 4}
	n4 := Node{Val: 4}
	n5 := Node{Val: 2}
	n6 := Node{Val: 1}
	n1.Left = &n2
	n1.Right = &n3
	n2.Left = &n4
	n2.Right = &n5
	n3.Right = &n6
}
