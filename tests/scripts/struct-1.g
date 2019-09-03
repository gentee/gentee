
struct s1 {
    int i
    str s
}

struct s2 {
    arr a2
    int j
}

struct s3 {
    map m3
    int k
    s2  ww
}

func getS s1 {
    s3 tmp = {k:33}
    s1 ret = {i: 7+tmp.k, s: "ok"}
    return ret
}