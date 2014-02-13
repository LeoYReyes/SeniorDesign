package main

import "fmt"

type Rectangle struct {
    length, width int
}

func (r *Rectangle) area() int {
    return r.length * r.width
}

func main() {
    rect := Rectangle{length: 5, width: 12}
    fmt.Println(rect.area())
}