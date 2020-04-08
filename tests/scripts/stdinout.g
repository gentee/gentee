run {
    while true {
        int val = Random(100)
        str num = ReadString("Enter \{val}:")
        Println(num)
        if num != str(val) : break
        sleep(1000)
    }
    Println("Test carriage")
    for i in 0..10 {
        Print("Percent: \{i*10}%\r")
        sleep(1000)
    }
}
