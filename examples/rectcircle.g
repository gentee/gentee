#!/usr/local/bin/gentee

// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

const : PI = 3.1415

func PrintResult(str text, float value) {
    Println("=================================
\{text}: \{Round(value, 4)}
=================================
")
}

func ReadFloat(str text) float {
    return float(ReadString(text + `: `))
}

run {
    bool ex
    while !ex {
        str action = ReadString(`Enter the number of the action:
1. Calculate the area of a rectangle
2. Calculate the area of a circle
3. Calculate the perimeter of a rectangle
4. Calculate the circumference of a circle
Press any other key to exit
`)
        switch action
        case `1`, `3` {
            float width = ReadFloat(`The width of the rectangle`)
            float height = ReadFloat(`The height of the rectangle`)
            if action == `1` {
                PrintResult("The area of the rectangle", width * height)
            } else {
                PrintResult("The perimeter of the rectangle", 2*(width + height))
            }
        }
        case `2` {
            float radius = ReadFloat(`The radius of the circle`)
            PrintResult("The area of the circle", PI * radius * radius)
        }
        case `4` {
            PrintResult("The circumference of a circle", 2 * PI * ReadFloat(`The radius of the circle`))
        }

        default: ex = true
    }
}
    