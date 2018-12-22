run str {
  str out = ReadString(`Enter your name: `)
  Println(`Hello, %{out}`)
  return out
}