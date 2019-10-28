#!/usr/local/bin/gentee
# result = 832040

/*
    http://rosettacode.org/wiki/Fibonacci_sequence

    The Fibonacci sequence is a sequence Fn of natural numbers defined recursively:

      F0 = 0 
      F1 = 1 
      Fn = Fn-1 + Fn-2, if n>1 

    Task
    Write a function to generate the nth Fibonacci number.

    Solutions can be iterative or recursive (though recursive solutions are generally considered too
    slow and are mostly used as an exercise in recursion).

    The sequence is sometimes extended into negative numbers by using a straightforward inverse of 
    the positive definition:

    Fn = Fn+2 - Fn+1, if n<0   
    support for negative n in the solution is optional.
*/

func fib(int n) int {
    if n==0 || n==1 : return n

    int prev=1
    int current=1
    for i in 2..n-1 {
        int next = prev + current
        prev = current
        current = next    
    }
    return current
}

run int {
    return fib(30)   
}
