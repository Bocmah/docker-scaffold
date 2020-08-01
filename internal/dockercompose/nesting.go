package dockercompose

import (
	"bufio"
	"strings"
)

type NestingLevel int

func (n NestingLevel) ApplyTo(str string) string {
	if n <= 0 {
		return str
	}

	spaces := strings.Repeat(" ", int(n)*2)
	scanner := bufio.NewScanner(strings.NewReader(str))

	var lines []string

	for scanner.Scan() {
		lines = append(lines, spaces+scanner.Text())
	}

	return strings.Join(lines, "\n")
}
