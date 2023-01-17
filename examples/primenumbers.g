#!/usr/local/bin/gentee
# stdin = 10000

// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

run {
    int high j
    str out
    set sieve

    || `This program uses "The Sieve of Eratosthenes" for finding prime numbers.
       `
    high = int( ReadString("Enter the high limit number ( < 100000 ): "))
    if high > 100000 : high = 100000

    for i in 2..high/2 {
        if !sieve[ i ] {
            j = i + i
            while j <= high {
               sieve[ j ] = true
               j += i
            }
        }
    }
    j = 0
    int width = *str(high) + 1

    for i in 2..high {
        if !sieve[ i ] {
            out += Format(`%%{width}d`,i)
            if ++j % 10 == 0 : out += "\r\n"
        }
    }
    str fname = JoinPath(TempDir(), "prime.txt")
    WriteFile( fname, out )
    || "\{j} prime numbers has been saved in \{fname}\n"
}
