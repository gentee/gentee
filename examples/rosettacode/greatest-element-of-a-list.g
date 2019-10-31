#!/usr/local/bin/gentee
# stdin = 4 56.7 -56.3 67.5 -33 -1000\n
# stdout = 1
# result = The greatest element is 67.5

/*
    http://rosettacode.org/wiki/Greatest_element_of_a_list

    Task
    Create a function that returns the maximum value in a provided set of values,
    where the number of values may not be known until run-time.
*/

run {
    arr nums = Split(ReadString(`Enter floats separated by spaces: `), ` `)
    if *nums > 0 {
        float fMax = float(nums[0])
        for fs in nums {
            float f = float(fs)
            if f > fMax : fMax = f
        }
        Println("The greatest element is \{fMax}")
    }
}

