package utils

import (
	"fmt"
	"strings"
)

func FormatMessage(counts map[string]int) string {
	builder := strings.Builder{}
	builder.WriteString("Code review rankings for today:\n")
	for user, count := range counts {
		builder.WriteString(fmt.Sprintf("%s: %d reviews\n", user, count))
	}
	return builder.String()
}
