#!/usr/local/bin/gentee
# stdin = 2020

// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

func  calendar( int year ) str {
    time  stime &= Date(year, 1, 1)

    str ret = `<!DOCTYPE html>
                <html lang="en"><head><meta charset="utf-8">
                <head><title>Calendar for year %{year}</title>
                <style type="text/css">
                    body {background: #FFF; font-family: Verdana;}
                    h1 {text-align: center; margin: 5px;}
                    table {border: 0; border-spacing: 7px;}
                    td {padding: 3px; border: 1px solid #777; text-align: center;}
                    #copy {font-size: smaller; text-align: center;}
                </style>
                </head>
                <body><center><h1>%{year}</h1>
                <table>`
    int firstday = 1
    int dayofweek = Weekday(stime)
    arr dayname
    str nline = "  \r\n"

    for i in 0..6 {
        dayname += Format( `ddd`, AddHours(stime, (i-dayofweek)*24))
    }

    for i in 0..3 {
        ret += `<tr>`
        for j in 1..3 {
            int month = i * 3 + j
            time first = Date(year, month, 1)
            ret += "<td>\{Format( `MMMM`, first )}<pre>"
            for k in 0..6 {
               ret += " " + dayname[(k+firstday) % 7]
            }
            ret += nline + Repeat( "    ", ( 7 + dayofweek - firstday ) % 7 )
 
            int lines
            for day in 1..Days( first ) {
                if dayofweek == 0 : ret += "<font color=red>"
                ret += Format("%4d", day)
                if dayofweek == 0 : ret += "</font>"
 
                dayofweek = ( dayofweek + 1 ) % 7
                if dayofweek == firstday {
                    ret += nline
                    lines++
                }
            }
            ret += Repeat("    ", ( 7 + firstday - dayofweek ) % 7 )
            while lines++ < 7 : ret += nline
            ret += `</pre></td>`
        }
        ret += `</tr>`
    }
    return ret + |`</table><br><div id="copy">Generated with Gentee Programming Language</div>
                   </body></html>`
}

run  {
    int year = int( ReadString( "Enter a year: "))
    str fname = JoinPath(TempDir(), `calendar.html`)
    WriteFile(fname, calendar(year))
    OpenWith(`firefox`, fname)
}