package diff

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type line struct {
	number int
	text   string
}

type deltaType int

const (
	Delete deltaType = iota
	Insert
	Equal
)

type delta struct {
	diffType deltaType
	oldLine  *line
	newLine  *line
}

func (d *delta) old() string {
	if d.oldLine != nil {
		return fmt.Sprintf("%2d", d.oldLine.number)
	}

	return "  "
}

func (d *delta) new() string {
	if d.newLine != nil {
		return fmt.Sprintf("%2d", d.newLine.number)
	}

	return "  "
}

func (d *delta) text() string {
	if d.oldLine != nil {
		return d.oldLine.text
	}

	return d.newLine.text
}

func Execute(args []string) {
	if len(args) != 2 {
		log.Fatalf("[ ERROR ] Incorrect number of arguments.\n" +
			"Usage is `sn diff [file1] [file2]`")
	}

	baseFile := args[0]
	diffFile := args[1]

	baseFileBytes, baseFileErr := os.ReadFile(baseFile)
	if baseFileErr != nil {
		log.Fatalf("[ ERROR ] Could not read file %s: %s\n", baseFile, baseFileErr.Error())
	}

	diffFileBytes, diffFileErr := os.ReadFile(diffFile)
	if diffFileErr != nil {
		log.Fatalf("[ ERROR ] Could not read file %s: %s\n", diffFile, diffFileErr.Error())
	}

	baseFileLines := lines(baseFileBytes)
	diffFileLines := lines(diffFileBytes)

	trace := buildTrace(baseFileLines, diffFileLines)
	diff := buildDiff(trace, baseFileLines, diffFileLines)
	outputDiff(diff)
}

func lines(bytes []byte) []line {
	lines := strings.Split(string(bytes), "\n")
	linesObj := make([]line, 0)
	for i, l := range lines {
		linesObj = append(linesObj, line{number: i + 1, text: l})
	}

	return linesObj
}

func buildTrace(baseFileLines, diffFileLines []line) [][]int {
	n := len(baseFileLines)
	m := len(diffFileLines)
	max := n + m

	v := make([]int, 2*max+1)
	trace := make([][]int, 0)

	for d := 0; d <= max; d++ {
		tempV := make([]int, len(v))
		copy(tempV, v)
		trace = append(trace, tempV)

		for k := -d; k <= d; k += 2 {
			var x int
			kIndex := k + max

			if k == -d || (k != d && v[kIndex-1] < v[kIndex+1]) {
				x = v[kIndex+1]
			} else {
				x = v[kIndex-1] + 1
			}

			y := x - k

			for x < n && y < m && baseFileLines[x].text == diffFileLines[y].text {
				x = x + 1
				y = y + 1
			}

			v[kIndex] = x

			if x >= n && y >= m {
				return trace
			}
		}
	}

	return trace
}

func buildDiff(trace [][]int, baseFileLines, diffFileLines []line) []delta {
	x := len(baseFileLines)
	y := len(diffFileLines)
	max := x + y

	diff := make([]delta, 0)
	for d := len(trace) - 1; d >= 0; d-- {
		v := trace[d]
		k := x - y

		var prevK int
		kIndex := k + max

		if k == -d || (k != d && v[kIndex-1] < v[kIndex+1]) {
			prevK = kIndex + 1
		} else {
			prevK = kIndex - 1
		}

		prevX := v[prevK]
		prevY := prevX - (prevK - max)

		for x > prevX && y > prevY {
			newDelta := createDelta(baseFileLines, diffFileLines, x-1, y-1, x, y)
			diff = append([]delta{newDelta}, diff...)
			x = x - 1
			y = y - 1
		}

		if d > 0 {
			newDelta := createDelta(baseFileLines, diffFileLines, prevX, prevY, x, y)
			diff = append([]delta{newDelta}, diff...)
		}

		x = prevX
		y = prevY
	}

	return diff
}

func createDelta(baseFileLines, diffFileLines []line, prevX, prevY, x, y int) delta {
	var aLine line
	var bLine line

	if prevX < len(baseFileLines) {
		aLine = baseFileLines[prevX]
	}

	if prevY < len(diffFileLines) {
		bLine = diffFileLines[prevY]
	}

	if x == prevX {
		return delta{diffType: Insert, oldLine: nil, newLine: &bLine}
	} else if y == prevY {
		return delta{diffType: Delete, oldLine: &aLine, newLine: nil}
	} else {
		return delta{diffType: Equal, oldLine: &aLine, newLine: &bLine}
	}
}

func outputDiff(diff []delta) {
	delete := "\033[48;2;89;64;67;31m"
	insert := "\033[48;2;64;89;67;32m"
	equals := "\033[39;0m"

	fileName := ".diffout"
	diffOut, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("[ ERROR ] Couldn't create diffout file: %s\n", err.Error())
	}

	for _, d := range diff {
		var colour string
		var tag string

		switch d.diffType {
		case Delete:
			colour = delete
			tag = "-"
		case Insert:
			colour = insert
			tag = "+"
		case Equal:
			colour = equals
			tag = " "
		default:
		}

		// log.Printf("%s%s %s %s    %s%s\n", colour, tag, d.old(), d.new(), d.text(), equals)
		diffString := fmt.Sprintf("%s%s %s %s    %s%s\n", colour, tag, d.old(), d.new(), d.text(), equals)
		_, err := diffOut.WriteString(diffString)
		if err != nil {
			log.Fatalf("[ ERROR ] Could not write to diffout file: %s\n", err.Error())
		}
	}

	if err := diffOut.Close(); err != nil {
		log.Fatal(err.Error())
	}

	cmd := exec.Command("less", "-f", "-R", "-F", "-X", fileName)
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		log.Fatal(err.Error())
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	if err := os.Remove(fileName); err != nil {
		log.Fatal(err.Error())
	}
}
