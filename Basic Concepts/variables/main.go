package main

import "fmt"

func main() {
	var name string = "Andrew"
	var age int = 30
	var height float64 = 1.75
	isActive := true

	fmt.Printf("Name: %s, Age: %d, Height: %.2f, Active: %t\n", name, age, height, isActive)
}
