//go:build !solution

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	data := os.Args[1:]

	overallStart := time.Now()

	log := make(chan string)

	for _, url := range data {
		go fetchURLInConcurrent(url, log)
	}

	for range data {
		fmt.Println(<-log)
	}

	overallFinish := time.Now()
	sec := overallFinish.Sub(overallStart).Seconds()

	fmt.Printf("%.2fs %s", sec, "overall time")
}

func fetchURLInConcurrent(url string, t chan string) {
	start := time.Now()

	response, err := http.Get(url)

	if err != nil {
		t <- fmt.Sprintf("Could not read %s", url)
		return
	}

	defer response.Body.Close()

	finish := time.Now()
	sec := finish.Sub(start).Seconds()

	t <- fmt.Sprintf("%.2fs %7d %s", sec, response.ContentLength, url)
}
