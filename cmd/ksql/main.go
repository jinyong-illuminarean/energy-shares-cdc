package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ksql <environment>")
		fmt.Println("Available environments: dev, qa, stage, prod")
		os.Exit(1)
	}

	fmt.Printf("KSQL Updater[%s]", os.Args[1])
}
