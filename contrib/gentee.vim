" gentee syntax highlighting for vim
"

if exists("b:current_syntax")
  finish
endif

syn region genteeString start=+L\="+ skip=+\\\\\|\\"+ end=+"+ contains=@Spell

syn keyword genteeStatement run

" Numerals
syn case ignore
"integer number, or floating point number without a dot and with "f".
syn match       genteeNumbers        display transparent "\<\d\|\.\d" contains=genteeNumber,genteeFloat,genteeOctError,genteeOct
syn match       genteeNumbersCom     display contained transparent "\<\d\|\.\d" contains=genteeNumber,genteeFloat,genteeOct
syn match       genteeNumber         display contained "\d\+\(u\=l\{0,2}\|ll\=u\)\>"

" hex number
syn match       genteeNumber         display contained "0x\x\+\(u\=l\{0,2}\|ll\=u\)\>"

" oct number
syn match       genteeOct            display contained "0\o\+\(u\=l\{0,2}\|ll\=u\)\>" contains=genteeOctZero
syn match       genteeOctZero        display contained "\<0"

syn match       genteeFloat          display contained "\d\+\.\d*\(e[-+]\=\d\+\)\="
syn match       genteeFloat          display contained "\d\+e[-+]\=\d\=\>"
syn match       genteeFloat          display "\(\.[0-9_]\+\)\(e[-+]\=[0-9_]\+\)\=[fl]\=i\=\>"

" Literals
syn region      genteeString         start=+L\="+ skip=+\\\\\|\\"+ end=+"+ contains=@Spell

syn match       genteeSpecial        display contained "\\\(x\x\+\|\o\{1,3}\|.\|$\)"
syn match       genteeCharacter      "L\='[^\\]'"
syn match       genteeCharacter      "L'[^']*'" contains=genteeSpecial


syn match       genteeFloat          display contained "\d\+\.\d*\(e[-+]\=\d\+\)\="
syn match       genteeFloat          display contained "\d\+e[-+]\=\d\=\>"
syn match       genteeFloat          display "\(\.[0-9_]\+\)\(e[-+]\=[0-9_]\+\)\=[fl]\=i\=\>"

syn keyword     genteeClause         import package
syn keyword     genteeConditional    if else switch
syn keyword     genteeBranch         goto break continue
syn keyword     genteeLabel          case default
syn keyword     genteeRepeat         for
syn keyword     genteeType      struct func
syn keyword     genteeType           int float bool str char arr map buf set obj handle

syn keyword     genteeTodo           contained TODO FIXME XXX
syn match       genteeLineComment    "\/\/.*" contains=@Spell,genteeTodo
syn match       genteeLineComment    "^#!.*$" contains=@Spell,genteeTodo

syn match       genteeCommentSkip    "^[ \t]*\*\($\|[ \t]\+\)"
syn region      genteeComment        start="/\*"  end="\*/" contains=@Spell,genteeTodo
syn region      genteeComment        start="||\`"  end="\`" contains=@Spell,genteeTodo
syn region      genteeComment        start="||\""  end="\"" contains=@Spell,genteeTodo


hi def link genteeStatement     Statement
hi def link genteeClause        Preproc
hi def link genteeConditional   Conditional
hi def link genteeBranch        Conditional
hi def link genteeLabel         Label
hi def link genteeRepeat        Repeat
hi def link genteeType          Type
hi def link genteeConcurrent    Statement
hi def link genteeValue         Constant
hi def link genteeBoolean       Boolean
hi def link genteeConstant      Constant
hi def link genteeBif           Function
hi def link genteeTodo          Todo
hi def link genteeLineComment   genteeComment
hi def link genteeComment       Comment
hi def link genteeNumbers       Number
hi def link genteeNumbersCom    Number
hi def link genteeNumber        Number
hi def link genteeFloat         Float
hi def link genteeOct           Number
hi def link genteeOctZero       Number
hi def link genteeString        String
hi def link genteeSpecial       Special
hi def link genteeCharacter     Character

let b:current_syntax = "gentee"
