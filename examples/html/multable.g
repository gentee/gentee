#!/usr/local/bin/gentee

// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

run  {
    str ret = `<!DOCTYPE html>
                <html lang="en"><head><meta charset="utf-8">
                <head><title>Multiplication Table</title>
                <style type="text/css">
                    body {background: #fff; font-family: Verdana;text-align: center;}
                    h1 {text-align: center; margin: 20px;}
                    table {border: 1px solid gray; border-collapse: separate;}
                    td,th {border: 1px solid gray;padding: 3px; width: 50px; text-align: center;}
                    th {background: #bbb;}
                    #copy {font-size: smaller; text-align: center;}
                </style>
                </head>
                <body><center><h1>Multiplication Table</h1>
                <table><tr><th>&nbsp;</th>`

    for i in 1..9 {
        ret += "<th>\{i}</th>"
    }   
    ret += `</tr>`
    for i in 1..9 {
        ret += "<tr><th>\{i}</th>"
        for j in 1..9 {
            ret += "<td>\{ i * j }</td>"
        }                        
        ret += `</tr>`
    }
    ret += `</table><br><div id="copy">Generated with Gentee Programming Language</div></body></html>`
    str fname = JoinPath(TempDir(), `multable.html`)
    WriteFile(fname, ret)
    OpenWith(`firefox`, fname)
}