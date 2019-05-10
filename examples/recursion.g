#!/usr/local/bin/gentee

// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

func factorial(int i) int {
    if i < 3 : return i
    return i * factorial(i-1)
}

func fibonacci( int pprev prev last) int {
    if last - 3 == 0 : return pprev + prev
    return fibonacci( prev, prev + pprev, last-1)
}

run {
    |`This program calculates 50th Fibonacci number and the factorial of 15
      Fibonacci number (Xn = Xn-1 + Xn-2)
      Factorial of n   (n! = 1*2*3*...*n)
     `
    int sum 
    int pprev = 1
    int prev = 1
    for i in 3..50 {
        sum = prev + pprev
        pprev = prev
        prev = sum
    }
    |`50th Fibonacci number = %{sum} (not recursive)
      50th Fibonacci number = %{fibonacci( 1, 1, 50)} (recursive)
    `

    int fact = 1
    for i in 1..15 : fact *= i
    |"15! = \{fact} (not recursive)
      15! = %{factorial(15)} (recursive)
     "
}
