#!/usr/local/bin/gentee
# stdout = 1
# result = CRC0x7ffe4151a6da6060

/*
    http://rosettacode.org/wiki/Loops/While

    Task
    Start an integer value at 1024.

    Loop while it is greater than zero.

    Print the value (with a newline) and divide it by two each time through the loop.
*/

run {
    int i = 1024
    while i > 0 {
        Println(i)
        i /= 2
    }
}