package tz

import (
	"testing"
)

func TestParseInput(t *testing.T) {
	successTests := [][]string{
		{"1234"},
		{"12:34pm"},
		{"12:34 pm"},
		{"0123"},
		{"1:23am"},
		{"1:23 am"},
		{"12:34pm 1/2/2000"},
		{"12:34 pm 01/02/2000"},
		{"1545 10/10/2023"},
		{"2301 1/1/2024"},
	}

	for _, test := range successTests {
		err := Execute(test)
		if err != nil {
			switch err.(type) {
			case *ParseError:
				t.Errorf("\nParse Error occurred: %s\nInvalid input: %s\n", err.Error(), test)

			default:
				t.Errorf("\nError occurred: %s\n", err.Error())
			}
		}
	}

	failTests := [][]string{
		{"2400"},
		{"123"},
		{"1260"},
		{"13:00pm"},
		{"13:00 am"},
		{"12:60am"},
		{"1:61 pm"},
		{"12:34am "},
		{"1234 32/12/2023"},
		{"1234 1/13/2023"},
		{"1234 001/12/2023"},
		{"1234 01/012/2023"},
		{"1234 01/12/02023"},
	}

	for _, test := range failTests {
		err := Execute(test)
		if err == nil {
			t.Errorf("\nExpected parse error for input: %s\n", test)
		}

		if _, ok := err.(*ParseError); !ok {
			t.Errorf("\nExpected parse error for input: %s\nGot: %s\n", test, err.Error())
		}
	}
}
