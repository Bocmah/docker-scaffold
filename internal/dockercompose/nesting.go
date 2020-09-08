package dockercompose

import (
	"bufio"
	"strings"
)

// NestingLevel represents how deep string to which it is applied should be nested
type NestingLevel int

// ApplyTo adds spaces (based on the NestingLevel) to the left of the string
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
