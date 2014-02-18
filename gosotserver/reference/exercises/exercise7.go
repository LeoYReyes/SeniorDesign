package main

import "fmt"
import "math"

type Rectangle struct {
    Positionable
    length, width int
}

type Positionable struct {
    x, y int
}

type Circle struct {
    Positionable
    radius float64
}

func (p *Positionable) isAtOrigin() bool {
    if p.x == 0 && p.y == 0 {
        return true
    }
    return false

}

func (c *Circle) area() float64 {
    return math.Pi * c.radius * c.radius
}
func (p *Positionable) area() int {
    return p.x * p.y
}
func (r *Rectangle) area() int {
    return r.length * r.width
}

func main() {
    rect := Rectangle{Positionable{10, 20}, 5, 12}
    fmt.Println(rect.area())
    circle := Circle{Positionable{0, 0}, 2}
    fmt.Println(rect.isAtOrigin())
    fmt.Println(circle.isAtOrigin())
    fmt.Println(circle.area())
}