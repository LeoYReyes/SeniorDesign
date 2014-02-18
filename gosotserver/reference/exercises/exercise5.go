package main

import "fmt"

func main() {

    xPtr, yPtr := new(int), new(int)
    
    fmt.Scan(xPtr)
    fmt.Scan(yPtr)
    
    fmt.Println(*xPtr + *yPtr)

}