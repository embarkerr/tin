package tz

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

type TZCommand struct {
	fs     *flag.FlagSet
	fromTZ string
	toTZ   string
	help   bool
}

type ParseError struct {
	msg string
}

func (e *ParseError) Error() string {
	return e.msg
}

func Execute(args []string) error {
	tzCommand := TZCommand{fs: flag.NewFlagSet("tzCommand", flag.ExitOnError)}
	tzCommand.fs.StringVar(&tzCommand.fromTZ, "f", "Local",
		"Provide the -f flag to specify timezone to convert from (default is Local)")
	tzCommand.fs.StringVar(&tzCommand.toTZ, "t", "UTC",
		"Provde -t to specify to specify timezone to convert to (default is UTC)")
	tzCommand.fs.BoolVar(&tzCommand.help, "h", false, "Help")
	tzCommand.fs.Parse(args)

	if tzCommand.help {
		printHelp()
		return nil
	}

	fromLocation, fromErr := time.LoadLocation(tzCommand.fromTZ)
	if fromErr != nil {
		log.Fatalf("[ ERROR ] Invalid -f location: %s\n", tzCommand.fromTZ)
	}

	toLocation, toErr := time.LoadLocation(tzCommand.toTZ)
	if toErr != nil {
		log.Fatalf("[ ERROR ] Invalid -t location: %s\n", tzCommand.toTZ)
	}

	localNow := time.Now()
	var timeRef = fmt.Sprintf("%02d%02d%02d", localNow.Hour(), localNow.Minute(), localNow.Second())
	var dateRef = fmt.Sprintf("%d/%d/%d", localNow.Day(), localNow.Month(), localNow.Year())
	var timeLayout = "150405 02/01/2006"
	if tzCommand.fs.NArg() > 0 {
		err := parseInput(&timeRef, &dateRef, &timeLayout, tzCommand.fs.Args())
		if err != nil {
			return err
		}
	}

	now, err := time.ParseInLocation(timeLayout, timeRef+" "+dateRef, fromLocation)

	if err != nil {
		log.Println("Parse error: " + err.Error())
	}

	locationTime := now.In(toLocation)
	fmt.Printf("Converting from: %s\n", now.String())
	fmt.Printf("Time in %s: %s\n", toLocation, locationTime)
	return nil
}

func parseInput(timeRef *string, dateRef *string, timeLayout *string, input []string) error {
	inputString := strings.Join(input, " ")

	// 12:34 pm or 12:34:56am
	timeOnlyRegex, err := regexp.Compile("^(?P<hour>(0?[0-9])|(1[0-2])):(?P<minute>[0-5][0-9])" +
		"(:(?P<second>[0-5][0-9]))?(\\s){0,1}(?P<period>AM|PM|am|pm)$")
	if err != nil {
		log.Fatalf(err.Error())
	}

	if timeOnlyRegex.MatchString(inputString) {
		template := "$hour,$minute,$second,$period"
		match := timeOnlyRegex.FindStringSubmatchIndex(inputString)
		var result []byte
		result = timeOnlyRegex.ExpandString(result, template, inputString, match)

		timeComponents := strings.Split(string(result), ",")
		hour := timeComponents[0]
		minute := timeComponents[1]
		second := timeComponents[2]
		period := timeComponents[3]

		if second == "" {
			second = "00"
		}

		*timeLayout = "03:04:05 pm 02/01/2006"
		*timeRef = fmt.Sprintf("%02s:%02s:%02s %02s", hour, minute, second, period)
		return nil
	}

	// 1234
	time24HourOnlyRegex, err := regexp.Compile("^(?P<hour>((0|1)[0-9])|(2[0-3]))(?P<minute>[0-5][0-9])" +
		"(?P<second>[0-5][0-9])?$")
	if err != nil {
		log.Fatalf(err.Error())
	}
	if time24HourOnlyRegex.MatchString(inputString) {
		template := "$hour,$minute,$second"
		match := time24HourOnlyRegex.FindStringSubmatchIndex(inputString)
		var result []byte
		result = time24HourOnlyRegex.ExpandString(result, template, inputString, match)

		timeComponents := strings.Split(string(result), ",")
		hour := timeComponents[0]
		minute := timeComponents[1]
		second := timeComponents[2]

		if second == "" {
			second = "00"
		}

		*timeLayout = "150405 02/01/2006"
		*timeRef = fmt.Sprintf("%02s%02s%02s", hour, minute, second)
		return nil
	}

	// 12:34 pm 18/1/2023
	dateTimeRegex, err := regexp.Compile("^(?P<hour>(0?[0-9])|(1[0-2])):(?P<minute>[0-5][0-9])" +
		"(:(?P<second>[0-5][0-9]))?(\\s){0,1}(?P<period>(AM|PM|am|pm)) " +
		"(?P<day>((0?[0-9])|((1|2)[0-9])|(3[0-1])))/(?P<month>((0?[0-9])|(1[0-2])))/(?P<year>\\d{0,4})$")
	if err != nil {
		log.Fatalf(err.Error())
	}
	if dateTimeRegex.MatchString(inputString) {
		template := "$hour,$minute,$second,$period,$day,$month,$year"
		match := dateTimeRegex.FindStringSubmatchIndex(inputString)
		var result []byte
		result = dateTimeRegex.ExpandString(result, template, inputString, match)

		dateTimeComponents := strings.Split(string(result), ",")
		hour := dateTimeComponents[0]
		minute := dateTimeComponents[1]
		second := dateTimeComponents[2]
		period := dateTimeComponents[3]
		day := dateTimeComponents[4]
		month := dateTimeComponents[5]
		year := dateTimeComponents[6]

		if second == "" {
			second = "00"
		}

		*timeLayout = "03:04:05 pm 02/01/2006"
		*timeRef = fmt.Sprintf("%02s:%02s:%02s %02s", hour, minute, second, period)
		*dateRef = fmt.Sprintf("%02s/%02s/%02s", day, month, year)
		return nil
	}

	// 1234 18/1/2023
	dateTime24HourRegex, err := regexp.Compile("^(?P<hour>((0|1)[0-9])|(2[0-3]))(?P<minute>[0-5][0-9])" +
		"(?P<second>[0-5][0-9])? " +
		"(?P<day>((0?[0-9])|((1|2)[0-9])|(3[0-1])))/(?P<month>((0?[0-9])|(1[0-2])))/(?P<year>\\d{0,4})$$")
	if err != nil {
		log.Fatalf(err.Error())
	}
	if dateTime24HourRegex.MatchString(inputString) {
		template := "$hour,$minute,$second,$day,$month,$year"
		match := dateTime24HourRegex.FindStringSubmatchIndex(inputString)
		var result []byte
		result = dateTime24HourRegex.ExpandString(result, template, inputString, match)

		dateTimeComponents := strings.Split(string(result), ",")
		hour := dateTimeComponents[0]
		minute := dateTimeComponents[1]
		second := dateTimeComponents[2]
		day := dateTimeComponents[3]
		month := dateTimeComponents[4]
		year := dateTimeComponents[5]

		if second == "" {
			second = "00"
		}

		*timeLayout = "150405 02/01/2006"
		*timeRef = fmt.Sprintf("%02s%02s%02s", hour, minute, second)
		*dateRef = fmt.Sprintf("%02s/%02s/%02s", day, month, year)
		return nil
	}

	log.Println("[ ERROR ] Invalid input format")
	return &ParseError{msg: "Invalid input format"}
}

func printHelp() {
	fmt.Printf("Timezone Converter\n" +
		"Usage:\n\t`sn tz [-flags] [input]`\n\n" +
		"[-flags]:\n" +
		"\t-f\ttimezone to convert from (default is Local)\n" +
		"\t-t\ttimezone to convert to (default is UTC)\n" +
		"\t-h\thelp\n\n" +
		"[input]:\n" +
		"\ttime to convert in 24 HR time format (e.g. 1830)\n")
}
