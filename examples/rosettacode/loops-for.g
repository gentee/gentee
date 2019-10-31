#!/usr/local/bin/gentee
# stdout = 1
# result = CRC0xcb9be5ca9477673f

/*
    http://rosettacode.org/wiki/Loops/For

    “For”   loops are used to make some block of code be iterated a number of times, setting a variable or parameter to a monotonically increasing integer value for each execution of the block of code.
    Common extensions of this allow other counting patterns or iterating over abstract structures other than the integers.

    Task
    Show how two loops may be nested within each other, with the number of iterations performed by the inner for loop being controlled by the outer for loop.

    Specifically print out the following pattern by using one for loop nested in another:
    *
    **
    ***
    ****
    *****
*/

run {
    for i in 0..4 {
        for j in 0..i : Print("*")
        Print("\n")
    }
}