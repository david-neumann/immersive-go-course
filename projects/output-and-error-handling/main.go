package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	url := "http://localhost:8080"
	maxRetries := 3

	for retries := 0; retries < maxRetries; retries++ {
		GetWeather(url)
	}

	fmt.Fprintf(os.Stderr, "Maximum retries reached. Unable to get a response from the server.\n")
}

func GetWeather(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to get a response from the server, please try again later.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", body)
		os.Exit(0)
	} else {
		retryAfterStr := resp.Header.Get("Retry-After")
		retryAfter, err := strconv.Atoi(retryAfterStr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to get a response from the server, please try again later.")
			os.Exit(1)
		}

		if retryAfter == 0 || retryAfter == 1 {
			time.Sleep(time.Duration(retryAfter) * time.Second)
		} else if retryAfter > 1 && retryAfter <= 5 {
			fmt.Fprintf(os.Stderr, "Server is taking longer than expected. Retrying after %d second(s)...\n", retryAfter)
			time.Sleep(time.Duration(retryAfter) * time.Second)
		} else {
			fmt.Fprintln(os.Stderr, "Server is taking too long to respond, please try again later.")
			os.Exit(1)
			return
		}
	}
}
