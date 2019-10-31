#!/usr/local/bin/gentee
# stdin = 50 -3\n
# stdout = 1
# result = CRC0x237f5e28886af9f6

/*
    http://rosettacode.org/wiki/Arithmetic/Integer

    Task
    Get two integers from the user, and then (for those two integers), display their:

    sum
    difference
    product
    integer quotient
    remainder
    exponentiation (if the operator exists)

    Don't include error handling.
    For quotient, indicate how it rounds   (e.g. towards zero, towards negative infinity, etc.).
    For remainder, indicate whether its sign matches the sign of the first operand or of the second operand, if they are different.
*/

run {
    arr nums = Split(ReadString(`Enter two integers separated by space or q to quit: `), ` `)
    if *nums == 2 {
        int a = int(nums[0])
        int b = int(nums[1])
        ||`%{a} + %{b} = %{a + b}
           %{a} - %{b} = %{a - b}
           %{a} * %{b} = %{a * b}
           %{a} / %{b} = %{a / b} // truncates towards 0
           %{a} % %{b} = %{a % b} // same sign as first operand
        `
    }
}

