# result = ok 777

run ok str {
    $GENTEE_Test = `ok %{777}`
    return $ %{$GOPATH}/bin/gentee scripts/env.g
}