run conout str {
    str ret = $ %{$GOPATH}/bin/gentee scripts/carriage.g
    buf bufout
    Run(`%{$GOPATH}/bin/gentee`, `scripts/carriage.g`, stdout: bufout)
    return ret + ClearCarriage(str(bufout))
}