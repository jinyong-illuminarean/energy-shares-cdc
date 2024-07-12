package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debezium-updater <environment>")
		fmt.Println("Available environments: dev, qa, stage, prod")
		os.Exit(1)
	}

	fmt.Printf("Debezium Updater[%s]", os.Args[1])
}
