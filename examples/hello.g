#!/usr/local/bin/gentee

// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

run {
    str name = ReadString(`Enter your name: `)
    Println(`Hello, %{ ?(*name>0, name, `world`) }!` )
}

/* Shorter versions
run : ||"Hello, world!\r\n"
run : $ echo "Hello, world!"
*/