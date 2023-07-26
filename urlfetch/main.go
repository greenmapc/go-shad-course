//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	input := os.Args[1:]

	for _, url := range input {
		response, err := http.Get(url)

		if err != nil {
			os.Exit(1)
		}

		body, _ := io.ReadAll(response.Body)
		fmt.Println(string(body))
		response.Body.Close()
	}

}
