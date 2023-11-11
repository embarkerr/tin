package b64

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Base64Command struct {
	fs          *flag.FlagSet
	decode      bool
	urlEncoding bool
	fileInput   string
	fileOutput  string
	help        bool
}

func Execute(args []string) string {
	// read flags
	base64Command := Base64Command{fs: flag.NewFlagSet("b64", flag.ExitOnError)}
	base64Command.fs.BoolVar(&base64Command.decode, "d", false, "Provide -d to decode")
	base64Command.fs.StringVar(&base64Command.fileInput, "f", "", "Provide the -f flag to read from input file")
	base64Command.fs.StringVar(&base64Command.fileOutput, "o", "", "Provide the -o flag to write the result to a file")
	base64Command.fs.BoolVar(&base64Command.urlEncoding, "u", false, "Provide the -u flag to use base64Url encoding")
	base64Command.fs.BoolVar(&base64Command.help, "h", false, "Help")
	base64Command.fs.Parse(args)

	if base64Command.help {
		printHelp()
		return ""
	}

	// parse input
	var input string
	if base64Command.fileInput != "" {
		data, err := os.ReadFile(base64Command.fileInput)
		if err != nil {
			log.Fatalf("[ ERROR ] Could not read file '%s': %s\n", base64Command.fileInput, err.Error())
		}
		input = string(data)
	} else if base64Command.fs.NArg() > 0 {
		input = strings.Join(base64Command.fs.Args(), " ")
	} else {
		fmt.Printf("[ ERROR ] Input is required to perform base64 conversions.\n" +
			"Usage is `sn [operation] [-flags] [input]`.\n")
		os.Exit(1)
	}

	// perform encode / decode
	var output string
	if base64Command.decode {
		var decode string
		var err error
		if base64Command.urlEncoding {
			decode, err = base64UrlDecode(input)
		} else {
			decode, err = base64Decode(input)
		}

		if err != nil {
			log.Fatalf("[ ERROR ] Could not decode '%s': %s\n", input, err.Error())
		}
		output = decode
	} else {
		var encode string
		if base64Command.urlEncoding {
			encode = base64UrlEncode(input)
		} else {
			encode = base64Encode(input)
		}
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

	return output
}

func base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func base64Decode(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	return string(data), err
}

func base64UrlEncode(input string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(input))
}

func base64UrlDecode(input string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(input)
	return string(data), err
}

func printHelp() {
	fmt.Printf("Base64 Encoder / Decoder\n" +
		"Usage:\n\t`sn b64 [-flags] [input]`\n\n" +
		"[-flags]:\n" +
		"\t-d\tdecode input (encode is default)\n" +
		"\t-f\tread input from file, provide file name to read from after flag\n" +
		"\t-o\twrite output to file, provide file name to write to after flag\n" +
		"\t-h\thelp\n\n" +
		"[input]:\n" +
		"\tstring input to encode/decode\n" +
		"\tif -f flag is provided, input becomes file to read input from\n" +
		"\tmake sure to escape special characters with `\\` where required.\n")
}
