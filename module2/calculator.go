package main

import (
    "errors"
    "fmt"
)

func Sum(a, b int) int {
    return a + b
}

func Subtract(a, b int) int {
    return a - b
}

func Multiply(a, b int) int {
    return a * b
}

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("can't divide by zero")
    }
    return a / b, nil
}

func calculator() {
    var a, b int
    var operation string

    fmt.Print("Enter the first number: ")
    fmt.Scanf("%d", &a)

    fmt.Print("Enter the second number: ")
    fmt.Scanf("%d", &b)

    fmt.Print("Enter the operation (add, subtract, multiply, divide): ")
    fmt.Scanf("%s", &operation)

    switch operation {
    case "add":
        fmt.Printf("The sum of the given numbers is: %d\n", Sum(a, b))
    case "subtract":
        fmt.Printf("The subtraction of the given numbers is: %d\n", Subtract(a, b))
    case "multiply":
        fmt.Printf("The product of the given numbers is: %d\n", Multiply(a, b))
    case "divide":
        if result, err := Divide(a, b); err != nil {
            fmt.Println(err)
        } else {
            fmt.Printf("The result of division is: %d\n", result)
        }
    default:
        fmt.Println("Invalid operation. Please enter a valid operation.")
    }
}

func main() {
    calculator()
}
