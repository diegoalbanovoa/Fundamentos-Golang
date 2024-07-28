package main

import "fmt"

func main() {
	a := 10
	b := 20

	sum := a + b
	difference := a - b
	product := a * b
	quotient := b / a
	remainder := b % a

	fmt.Printf("Sum: %d, Difference: %d, Product: %d, Quotient: %d, Remainder: %d\n", sum, difference, product, quotient, remainder)
}
