
run cmdline str {
    str out
    arr args = Args()
    switch ArgCount() 
    case 0 : out = str(IsArg(`-`)) +  str(IsArg(``)) + Join(ArgTail(), `=`) + Arg(`-n`, `ok`)
    case 1 : out = str(ArgCount()) + args[0]
    case 2 {
        out = Arg(`p`) + Arg(`--flag`, `empty`) + Arg(`-n`, `false`)
    } 
    case 3 : out = Join(ArgList(`-list`), `+`)
    case 4 : out = Arg(`o`) + `+` + Join(ArgTail(), `+`)
    case 5 : out = Join(ArgTail(), `+`) + str(IsArg(`i`)) + str(IsArg(`one`)) + str(ArgInt(`i`, 99))
    case 6 : out = Join(ArgTail(), `+`) + str(IsArg(`-`)) +  str(IsArg(``))
    return out
}