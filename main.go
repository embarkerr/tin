package main

import (
	"fmt"
	"log"
	"os"
	"sn/commands/b64"
	"sn/commands/diff"
)

func main() {
	args := os.Args
	if len(args) > 1 && args[1] == "-h" {
		fmt.Printf("\t== Tin v0.1 ==\n" +
			"Usage:\n\t`sn [command] [options]`\n\n" +
			"Commands:\n" +
			"\tb64\tBase64 encoder / decoder\n" +
			"\tdiff\tFile diff tool\n" +
			"\tmore to come\n\n" +
			"Use `sn [command] -h` for help on different commands\n")
		return
	}

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
