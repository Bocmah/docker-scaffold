package assemble

import (
	"strings"
)

func formatAppName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}
