#!/usr/local/bin/gentee
# settings.cycle = 25
# result = .../examples/rosettacode/loops-infinite.g [13:5] maximum cycle count has been reached

/*
    http://rosettacode.org/wiki/Loops/Infinite

    Task
    Print out SPAM followed by a newline in an infinite loop.
*/

run {
    while true : Println(`SPAM`)
}