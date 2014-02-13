package main

import ("fmt"; "time")

func send(c chan int) {
    for i := 10; i >= 0; i-- {
        c <- i
    }
}

func main() {
    c := make(chan int)
    go send(c)
    for <-c > 0 {
        fmt.Println(<-c)
        time.Sleep(time.Duration(1) * time.Second)
    }

}