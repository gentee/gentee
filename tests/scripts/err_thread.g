func divZero() int {
    return 10/0
}

func errThread {
    go {
        int i = divZero()
    }
    Sleep(3000)
}

func errThread2 {
    go {
        error(1000, `This is an error message`)
        int i = divZero()
    }
    Sleep(3000)
}