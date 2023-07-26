//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	
	fileNames := os.Args[1:]
	counterMap := make(map[string]int)

	for _, fileName := range(fileNames) {
		data, err := os.ReadFile(fileName)
		
		check(err)
		lines := strings.Split((string(data)), "\n")

		for _, line := range lines {
			existingValue, exists := counterMap[string(line)]
			
			if exists {
				counterMap[string(line)] = existingValue + 1
			} else {
				counterMap[string(line)] = 1
			}
		}
	}

	for key, value := range counterMap {
		if value >= 2 {
			fmt.Print(value, "\t", key)
		}
		fmt.Println()
	}
}
