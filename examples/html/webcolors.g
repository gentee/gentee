#!/usr/local/bin/gentee

// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

const {
    LCOUNT = 12
}
 
func item( str rgb ) str {
   return `<td><div style="background-color: #%{rgb}">&nbsp;</div>%{rgb}</td>`
}

run  {
    str ret = `<!DOCTYPE html>
                <html lang="en"><head><meta charset="utf-8">
                <head><title>RGB colors for Web</title>
                <style type="text/css">
                    body {background: #fff; font-family: Verdana;text-align: center;}
                    h1 {text-align: center; margin: 20px;}
                    td {padding: 5px; width: 50px; text-align: center;}
                    td div {padding: 5px;}
                    #copy {font-size: smaller; text-align: center;}
                </style>
                </head>
                <body><center><h1>RGB colors for Web</h1>
                <table><tr>`

    int cur
    
    local outitem(int rgb )  {
        ret += item( Format( "%06X", rgb ))
        if ++cur == LCOUNT {
            ret += `</TR><TR>`
            cur = 0 
        }
    }
    int i = 0xFF
    while i >= 0 {
        int j = 0xFF
        while j >= 0 {
            int k = 0xFF
            while k >= 0 {
               outitem(( i << 16 ) + ( j << 8 ) + k)
               k -= 0x33
            }     
            j -= 0x33
        }     
        i -= 0x33
    }
    local colors( int start step ) {
        while start > 0 {
            outitem(start)
            start -= step
        }
    }
    colors(0xFFFFFF, 0x111111)
    outitem(0)
    colors(0xFF0000, 0x110000)
    colors(0x00FF00, 0x001100)
    colors(0x0000FF, 0x000011)

    ret += `</table><br><div id="copy">Generated with Gentee Programming Language</div></body></html>`
    str fname = JoinPath(TempDir(), `webcolors.html`)
    WriteFile(fname, ret)
    OpenWith(`firefox`, fname)
}