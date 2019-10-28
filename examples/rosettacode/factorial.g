#!/usr/local/bin/gentee
# result = 1307674368000

/*
    http://rosettacode.org/wiki/Factorial

    Definitions
    The factorial of 0 (zero) is defined as being 1 (unity).
    The Factorial Function of a positive integer, n, is defined as the product of the sequence:
    n,   n-1,   n-2,   ...   1 

    Task
    Write a function to return the factorial of a number.
    Solutions can be iterative or recursive.
    Support for trapping negative n errors is optional.
*/

func factorial(int n) int {
    if n < 0 : error(1000, "argument less than 0")
    int fact = 1
    for i in 1..n : fact *= i
    return fact
}
 
run int {
    return factorial(15)
}