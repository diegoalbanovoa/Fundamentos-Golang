# Project Week 1: Go Fundamentals

## Description

This project covers the basic concepts of the Go language, including basic syntax, control structures, variables and data types, and creating a simple Go program.

## Project Structure

The project is divided into several subdirectories, each containing specific examples:

- **hello_world**: Basic "Hello, World" program.
- **variables**: Examples of variables and data types.
- **constants**: Usage of constants in Go.
- **operators**: Examples of operators and expressions.
- **if_else**: Conditional control structures.
- **for_loop**: Examples of `for` loops.

## Content

### 1. Hello World
File: `hello_world/main.go`

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```
## 2. Variables and Data Types
File: `variables/main.go`

```go
package main

import "fmt"

// main is the entry point of the program
func main() {
    // Define variables with different data types
    var name string = "Andrew"
    var age int = 30
    var height float64 = 1.75
    isActive := true // short variable declaration

    // Print the values of the variables
    fmt.Printf("Name: %s, Age: %d, Height: %.2f, Active: %t\n", name, age, height, isActive)
}
```
## 3. Constants
File: `constants/main.go`

```go
package main

import "fmt"

// Define constants
const Pi = 3.14
const E = 2.71

// main is the entry point of the program
func main() {
    // Print the values of the constants
    fmt.Println("Pi:", Pi)
    fmt.Println("E:", E)
}
```

## 4. Operators and Expressions
File: `operators/main.go`

```go
package main

import "fmt"

// main is the entry point of the program
func main() {
    a := 10
    b := 20

    // Perform basic arithmetic operations
    sum := a + b
    difference := a - b
    product := a * b
    quotient := b / a
    remainder := b % a

    // Print the results of the operations
    fmt.Printf("Sum: %d, Difference: %d, Product: %d, Quotient: %d, Remainder: %d\n", sum, difference, product, quotient, remainder)
}
```
## 5. Control Structures: If and Else
File: `if_else/main.go`

```go
package main

import "fmt"

// main is the entry point of the program
func main() {
    age := 20

    // Use an if-else statement to check the age
    if age >= 18 {
        fmt.Println("You are an adult.")
    } else {
        fmt.Println("You are a minor.")
    }
}
```
## 6. Control Structures: For Loop
File: `for_loop/main.go`

```go
package main

import "fmt"

// main is the entry point of the program
func main() {
    // Use a for loop to iterate five times
    for i := 0; i < 5; i++ {
        fmt.Println("Iteration:", i)
    }
}
```

