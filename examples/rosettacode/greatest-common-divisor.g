#!/usr/local/bin/gentee
# stdin = 49865 69811\n
# stdout = 1
# result = GCD is 9973

/*
    http://rosettacode.org/wiki/Greatest_common_divisor

    Task
    Find the greatest common divisor of two integers.
*/

func gcd( int left right ) int {
    if right == 0 : return left
    return gcd( right, left % right )
}

run {
    arr nums = Split(ReadString(`Enter two inetegers separated by space: `), ` `)
    if *nums == 2 {
        Println("GCD is \{gcd(int(nums[0]), int(nums[1]))}")
    }
}

