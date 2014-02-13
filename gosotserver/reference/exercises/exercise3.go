package main

import "fmt"

func main() {
	romanMap := make(map[string]int)
	romanMap["I"] = 1
	romanMap["II"] = 2
	romanMap["III"] = 3
	romanMap["IV"] = 4
	romanMap["V"] = 5
	slice := make([]string, 2, 2)
	for i := 0; i < 2; i++ {
		fmt.Scan(&slice[i])
	}
	for i := 0; i < 2; i++ {
		if _, found := romanMap[slice[i]]; !found {
			fmt.Println("Invalid input")
		}
	}
	fmt.Println(romanMap[slice[0]] + romanMap[slice[1]])

}
