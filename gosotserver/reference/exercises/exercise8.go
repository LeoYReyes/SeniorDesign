package main

import "fmt"

type Shaper interface {
    Area() int
    print()
}
type Rectangle struct {
    length, width int
}
type Square struct {
    side int
}
func (r Rectangle) Area() int {
    return r.length * r.width
}
func (s Square) Area() int {
    return s.side * s.side
}
func (s Square) print() {
    fmt.Println(s.Area())
}
func (r Rectangle) print() {
    fmt.Println(r.Area())
}

func main() {
    r := Rectangle{2, 3}
    s := Square{5}
    s1, s2 := Shaper(r), Shaper(s)
    
    s1.print()
    s2.print()
}