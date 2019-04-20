include {
    "err_thread.g"
}

run int {
    go : errThread()
    Sleep(3000)
    return 0
}