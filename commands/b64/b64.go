package b64

import (
	"encoding/base64"
	"flag"
	"log"
	"os"
)

type Base64Command struct {
	fs         *flag.FlagSet
	decode     bool
	fileInput  string
	fileOutput string
}

func Execute(args []string) {
	// read flags
	base64Command := Base64Command{fs: flag.NewFlagSet("b64", flag.ExitOnError)}
	base64Command.fs.BoolVar(&base64Command.decode, "d", false, "Provide -d to decode")
	base64Command.fs.StringVar(&base64Command.fileInput, "f", "", "Provide the -f flag to read from input file")
	base64Command.fs.StringVar(&base64Command.fileOutput, "o", "", "Provide the -o flag to write the result to a file")
	base64Command.fs.Parse(args)

	// parse input
	var input string
	if base64Command.fileInput != "" {
		data, err := os.ReadFile(base64Command.fileInput)
		if err != nil {
			log.Fatalf("[ ERROR ] Could not read file '%s': %s\n", base64Command.fileInput, err.Error())
		}
		input = string(data)
	} else if base64Command.fs.NArg() == 1 {
		input = base64Command.fs.Args()[0]
	} else {
		log.Fatalf("[ ERROR ] Input is required to perform base64 conversions.\n" +
			"Usage is `sn [operation] [-flags] [input]`.\n")
	}

	// perform encode / decode
	var output string
	if base64Command.decode {
		decode, err := base64Decode(input)
		if err != nil {
			log.Fatalf("[ ERROR ] Could not decode '%s': %s\n", input, err.Error())
		}
		output = decode
	} else {
		encode := base64Encode(input)
		output = encode
	}

	// send output
	if base64Command.fileOutput != "" {
		err := os.WriteFile(base64Command.fileOutput, []byte(output), os.ModePerm)
		if err != nil {
			log.Fatalf("[ ERROR ] Failed to write to file %s\nOutput was: %s\n", base64Command.fileOutput, output)
		}
	} else {
		log.Printf("Output:\n%s\n", output)
	}
}

func base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func base64Decode(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	return string(data), err
}
