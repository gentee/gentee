#!/usr/local/bin/gentee
# stdout = 1
# result = Len:6 [7 5 4 6 2 1]

/*
    http://rosettacode.org/wiki/Arrays

    This task is about arrays.

    For hashes or associative arrays, please see Creating an Associative Array.
    For a definition and in-depth discussion of what an array is, see Array.

    Task
    Show basic array syntax in your language.

    Basically, create an array, assign a value to it, and retrieve an element (if available, 
    show both fixed-length arrays and dynamic arrays, pushing a value into it).
*/

run {
    arr.int ai = {1,2,3,4,5}
    ai[2] = ai[1] + ai[3]
    ai += 7
    Println("Len:\{*ai}", Reverse(ai))
}