run str {
  str out = ReadString(`Enter your name: `)
  Println(out = `Hello, %{out}`)
  str fname = JoinPath(TempDir(), `hello.txt`)
  WriteFile(fname, out)
  Open(fname)
  return out
}