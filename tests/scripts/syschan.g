func mythread(int id) {
    for i in 1..10 {
        Print(" \{id}=\{i} ")
        sleep(1000 + id*200)
    }
}

run syschan {
    for i in 1..3 {
        go (id: i) : mythread(id)
        sleep(100)
    }
    for i in 1..20 {
        Print("  \{0}=\{i} ")
        sleep(1000)
    }
}