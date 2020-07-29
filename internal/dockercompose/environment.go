package dockercompose

import (
	"fmt"
	"strings"
)

type Environment map[string]string

func (e Environment) Render() string {
	length := len(e)

	if length == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("environment:\n")

	i := 1
	for variable, value := range e {
		if variable == "" {
			continue
		}

		if value == "" {
			sb.WriteString(fmt.Sprintf("  %s:", variable))
		} else {
			sb.WriteString(fmt.Sprintf("  %s: %s", variable, value))
		}

		if i != length {
			sb.WriteString("\n")
		}

		i++
	}

	return sb.String()
}
