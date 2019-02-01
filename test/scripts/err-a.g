func c_func(int a) int {
    return a*2
}

include {
    "cΣ.g"
}

run int {
    return c_func(10)
}