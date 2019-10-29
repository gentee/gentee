#!/usr/local/bin/gentee
# stdin = -99 200\nq
# stdout = 1
# result = 101

/*
    http://rosettacode.org/wiki/A%2BB

    A+B - a classic problem in programming contests, it's given so contestants can gain familiarity
    with the online judging system being used.

    Task
    Given two integers, A and B.
    Their sum needs to be calculated.

    Input data
    Two integers are written in the input stream, separated by space(s):
    -1000 <= A,B <= +1000

    Output data
    The required output is one integer: the sum of A and B.
*/

run {
    while true {
        arr nums = Split(ReadString(`Enter two integers separated by space or q to quit: `), ` `)
        if *nums == 1 && nums[0] == `q` : break
        if *nums != 2 {
            Println(`Invalid input, try again`)
            continue
        } 
        int a = int(nums[0])
        int b = int(nums[1])
        if Abs(a) > 1000 || Abs(b) > 1000 {
            Println("Both numbers must be in the interval [-1000, 1000], try again")
            continue
        }
        Println(a+b)
    }
}