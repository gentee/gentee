func myerr(int num) int {
    while num < 20 {
        if num == 10 : error( 30*num, `Σ custom error №%d`, num/2 )
        num++
    }
    return 0
}
run int {
    return myerr(1)
}