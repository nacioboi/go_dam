package main

import "fmt"

func _asm_add(a, b int) int

func main_asm() {
	x := _asm_add(1, 2)
	fmt.Println(x)
}
