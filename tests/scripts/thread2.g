include {
    "err_thread.g"
}

run int {
    go : errThread2()
    return 0
}