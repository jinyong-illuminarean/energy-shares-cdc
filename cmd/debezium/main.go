package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debezium-updater <environment>")
		fmt.Println("Available environments: dev, qa, stage, prod")
		os.Exit(1)
	}

	fmt.Printf("Debezium Updater[%s]\n", os.Args[1])

	response, err := http.Get("https://dbzm-common-gw.illuminarean.com")
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	fmt.Println("Status Code:", response.StatusCode)
	fmt.Println("Content Length:", response.ContentLength)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	fmt.Println("Response Body:")
	fmt.Println(string(body))
}
