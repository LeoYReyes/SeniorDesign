package main

import "fmt"

func main() {
	slice := make([]string, 5, 5)
	for i := 0; i < 5; i++ {
		fmt.Scan(&slice[i])
	}
	slice = append(slice, slice...)
	for key := range slice {
		fmt.Println(slice[key])
	}
}
