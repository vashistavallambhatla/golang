package main

import (
    "fmt"
    "errors"
)

func Sum(a int,b int) int {
    return a+b
}
func Subtract(a int,b int) int {
    return a-b
}
func Multiply(a int,b int) int {
    return a*b
}
func Divide(a int,b int) (int,error) {
    if b == 0 {
        return 0, errors.New("Can't divide by zero")
    }
    return a/b , nil
}
func calculator() {
    fmt.Println("Enter the first number:")
    var a int
    fmt.Scanf("%d",&a)
    fmt.Println("Enter the second number:")
    var b int 
    fmt.Scanf("%d",&b)
    fmt.Println("Enter the operation:")
    var operation string
    fmt.Scanf("%s",&operation)
    
    if operation == "addition" {
        fmt.Println("The sum of the given numbers is : %d\n",Sum(a,b))
    } else if operation == "subtraction" {
        fmt.Println("The subtraction of the given numbers is : %d\n",Subtract(a,b))
    } else if operation == "multiplication" {
        fmt.Println("The product of the given numbers is : %d\n",Multiply(a,b))
    } else if operation == "division" {
       result,err := Divide(a,b)
        if err != nil {
            fmt.Println(err)
        }else {
            fmt.Println("The result of division : %d\n",result)
        }
    }
}
