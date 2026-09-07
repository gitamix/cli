// main defines the main package for the project.
package main

import (
	"log"

	"github.com/gitamix/cli/cmd/gitamix"
)

// main is the entry point for the command.
func main() {
	err := gitamix.Execute()
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}
}
