
run cmdline str {
    str out
    arr args = Args()
    switch ArgCount() 
    case 0 : out = str(IsArg(`-`)) +  str(IsArg(``)) + Join(ArgsTail(), `=`) + Arg(`-n`, `ok`)
    case 1 : out = str(ArgCount()) + args[0]
    case 2 {
        out = Arg(`p`) + Arg(`--flag`, `empty`) + Arg(`-n`, `false`)
    } 
    case 3 : out = Join(Args(`-list`), `+`)
    case 4 : out = Arg(`o`) + `+` + Join(ArgsTail(), `+`)
    case 5 {
       out = Join(ArgsTail(), `+`) + str(IsArg(`i`)) + str(IsArg(`one`)) + str(Arg(`i`, 99)) + 
             str(IsArg(`two`))
    }
    case 6 : out = Join(ArgsTail(), `+`) + str(IsArg(`-`)) +  str(IsArg(``))
    return out
}