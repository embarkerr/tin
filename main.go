package main

import (
	"log"
	"os"
	"sn/commands/b64"
	"sn/commands/diff"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("[ ERROR ] Not enough arguments")
	}

	command := args[1]
	switch command {
	case "b64":
		b64.Execute(args[2:])

	case "diff":
		diff.Execute(args[2:])

	default:
		log.Fatalf("[ ERROR ] Command '%s' is not supported\n", command)
	}
}
