package main

import (
	"fmt"
	"reflect"
)

type One struct {
	Val string
}

func (o One) get() {
	fmt.Println(reflect.TypeOf(o))
	o.Val = "empty"
}

func main() {
	p := One{"H"}
	fmt.Println(reflect.TypeOf(p))
	(&p).get()
	fmt.Println(p)
}
