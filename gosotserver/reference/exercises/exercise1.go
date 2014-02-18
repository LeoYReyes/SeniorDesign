package main

import "fmt"

func main() {
	i := 1
	for i <= 10 {
		if i % 4 == 0 {
			fmt.Println(i, ":-)")
		} else {
			fmt.Println(i, ":-(")
		}
		i = i + 1
	}
}
