package main

import (
	"errors"
	"fmt"
	"math"
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
		return nil, errors.New("Can't pop from the empty stack")
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

func (st *Stack) PopLazy() (*Node, error) {
	if len(st.Nodes) <= 0 {
		return nil, errors.New("Can't pop from the empty stack")
	}
	n := st.Nodes[len(st.Nodes)-1]
	st.Nodes = st.Nodes[:len(st.Nodes)-1]
	return n, nil
}

func depthFirstSearch(root *Node) []interface{} {
	results := make([]interface{}, 0)
	if root == nil || root.Val == nil {
		return results
	}

	stack := Stack{}
	stack.Push(root)
	for len(stack.Nodes) != 0 {
		n, _ := stack.Pop()
		results = append(results, n.Val)
	}
	return results
}

func depthFirstSearch_rec(root *Node) []interface{} {
	results := make([]interface{}, 0)
	if root == nil || root.Val == nil {
		return results
	}
	results = append(results, root.Val)
	results = append(results, depthFirstSearch_rec(root.Left)...)
	results = append(results, depthFirstSearch_rec(root.Right)...)
	return results
}

func treeIncludes(root *Node, target interface{}) bool {
	stack := Stack{}
	stack.Push(root)
	//       a
	//     /   \
	//   b      c
	//  / \      \
	// d   e	  f
	for len(stack.Nodes) > 0 {
		current, _ := stack.Pop()
		if current.Val == target {
			return true
		}
		if current.Right != nil {
			stack.Push(current.Right)
		}
		if current.Left != nil {
			stack.Push(current.Left)
		}
	}

	return false
}

func treeIncludes_rec(root *Node, target interface{}) bool {
	if root == nil || root.Val == nil {
		return false
	}
	fmt.Printf("%v -> ", root.Val)
	if root.Val == target {
		return true
	}

	return treeIncludes_rec(root.Left, target) || treeIncludes_rec(root.Right, target)
}

func treeSum_rec(root *Node) int {
	sum := 0
	if root == nil || root.Val == nil {
		return sum
	}
	sum += root.Val.(int) + treeSum_rec(root.Left) + treeSum_rec(root.Right)
	return sum
}

func treeMin(root *Node) int {
	min := math.MaxInt64
	if root == nil || root.Val == nil {
		return min
	}

	stack := Stack{}
	stack.Push(root)
	current, _ := stack.PopLazy()
	for ; current != nil; current, _ = stack.PopLazy() {
		cmin, ok := current.Val.(int)
		if ok && cmin < min {
			min = cmin
		}
		if current.Left != nil {
			stack.Push(current.Left)
		}
		if current.Right != nil {
			stack.Push(current.Right)
		}
	}
	return min
}

func treeMin_rec(root *Node) int {
	min := math.MaxInt64
	if root == nil || root.Val == nil {
		return min
	}

	cmin, ok := root.Val.(int)
	if !ok {
		return min
	}
	min = cmin
	lmin := treeMin_rec(root.Left)
	if lmin < min {
		min = lmin
	}
	rmin := treeMin_rec(root.Right)
	if rmin < min {
		min = rmin
	}
	return min
}

func maxRootToLeafSum(root *Node) int {
	max := math.MinInt64
	if root == nil || root.Val == nil {
		return max
	}
	max = root.Val.(int)
	if root.Left == nil && root.Right == nil {
		return max
	}
	lmax := maxRootToLeafSum(root.Left)
	rmax := maxRootToLeafSum(root.Right)
	if lmax > rmax {
		max += lmax
	}
	if rmax >= lmax {
		max += rmax
	}
	return max
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
	res := depthFirstSearch(&a)
	fmt.Println(res)
	rec_res := depthFirstSearch_rec(&a)
	fmt.Println(rec_res)
	fmt.Println(treeIncludes_rec(&a, "c"))
	fmt.Println(treeIncludes_rec(&a, "g"))
	fmt.Println(treeIncludes(&a, "c"))
	fmt.Println(treeIncludes(&a, "k"))
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
	// fmt.Println(treeSum_rec(&n1))
	// fmt.Println(treeSum_rec(&n1))
	// fmt.Println(treeMin_rec(&n1))
	// fmt.Println(treeMin(&n1))
	// fmt.Println(maxRootToLeafSum(&n1))
	// fmt.Println(maxRootToLeafSum(nil))
	fmt.Println(maxRootToLeafSum(&Node{Val: 10, Left: nil, Right: nil}))
}
