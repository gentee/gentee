#result = s0[company:My company name:My name owner:s1[i:40 s:ok] b:sq[a:54] f:s2[a2:[1 2] j:6]]

include : "struct-1.g"

struct sq {
    int a
}

struct s0 {
    str company
    str name
    s1  owner
    sq  b
    s2  f
}

run q s0 {
   s0 ret
   s1 t = {s: `as`, i:45}
   sq q = {a: 9+t.i}
   ret.company = "My company"
   ret.owner = getS()
   ret.name = "My name"
   ret.b &= q
   s1 f = t
   s2 f2 = {j: 6, a2:{`1`, `2`}}
   ret.f = f2
   return ret
}
