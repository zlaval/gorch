package main

import (
	"gorch/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("execution error: %v", err)
	}
}
