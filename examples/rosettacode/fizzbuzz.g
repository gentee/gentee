#!/usr/local/bin/gentee
# result = CRC0x528c281d97aef54e
# stdout = 1

/*
    http://rosettacode.org/wiki/FizzBuzz

    Task
    Write a program that prints the integers from 1 to 100 (inclusive).

    But:
        for multiples of three, print Fizz (instead of the number)
        for multiples of five, print Buzz (instead of the number)
        for multiples of both three and five, print FizzBuzz (instead of the number)

    The FizzBuzz problem was presented as the lowest level of comprehension required to illustrate
    adequacy.
*/

run {
    for i in 1..100 {
        str out
        if i % 3 == 0 : out = "Fizz"
        if i % 5 == 0 : out += "Buzz"
        if *out == 0 : out = str(i)
        Println(out)
    } 
}