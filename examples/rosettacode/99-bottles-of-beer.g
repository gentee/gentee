#!/usr/local/bin/gentee
# result = CRC0xc1cef31505cd8c50
# stdout = 1

/*
    http://rosettacode.org/wiki/99_Bottles_of_Beer

    Task
    Display the complete lyrics for the song: 99 Bottles of Beer on the Wall.

    The beer song
    The lyrics follow this form:

    99 bottles of beer on the wall
    99 bottles of beer
    Take one down, pass it around
    98 bottles of beer on the wall

    98 bottles of beer on the wall
    98 bottles of beer
    Take one down, pass it around
    97 bottles of beer on the wall

    ... and so on, until reaching 0.

    Grammatical support for "1 bottle of beer" is optional.

    As with any puzzle, try to do it in as creative/concise/comical a way as possible (simple, obvious solutions allowed, too).
*/

run {
   	local bottles(int i) str {
		if i==0: return "No more bottles"
		if i==1: return "1 bottle"
		return "\{i} bottles"
	}
 
	for i in 99..1 {
        str s = bottles(i)
        ||`%{s} of beer on the wall
           %{s} of beer
		   Take one down, pass it around
		   %{bottles(i-1)} of beer on the wall
        
        `
	}
}