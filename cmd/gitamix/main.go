package main

import (
	"log"
)

// main is the entry point for the command.
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("panic recovered: %v", r)
		}
	}()
	log.Print("it works")
}
