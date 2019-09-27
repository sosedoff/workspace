package main

import (
	"fmt"
	"log"
	"os"
)

func exitWithError(err interface{}) {
	fmt.Println("error:", err)
	os.Exit(1)
}

func requireConfirmation(message, accepted string) bool {
	var answer string
	fmt.Printf("%s (yes/no): ", message)

	_, err := fmt.Scanln(&answer)
	if err != nil {
		log.Fatal(err)
	}

	return answer == "yes"
}
