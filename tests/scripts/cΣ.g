func c_func(int i) int : return i + 5

pub func cpub_func(int i) int : return i*3

run int {
    return cpub_func(10) + c_func(5)
}
