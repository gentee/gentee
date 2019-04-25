include {
    "err_thread.g"
}

run int {
    go : errThread()
    sleep(3000)
    return 0
}