package utils

import (
	"fmt"
	"sort"
	"strings"
)

type keyValue struct {
	Key   string
	Value int
}

func FormatMessage(counts map[string]int) string {
	if len(counts) == 0 {
		return "No reviews today :("
	}

	var sortedSlice []keyValue
	for key, value := range counts {
		sortedSlice = append(sortedSlice, keyValue{key, value})
	}

	// Actually sort the slice by value in descending order
	sort.Slice(sortedSlice, func(i, j int) bool {
		return sortedSlice[i].Value > sortedSlice[j].Value
	})

	var builder strings.Builder
	builder.WriteString("Code review rankings for today:\n")
	for _, ranking := range sortedSlice {
		builder.WriteString(fmt.Sprintf("%s: %d reviews\n", ranking.Key, ranking.Value))
	}

	return builder.String()
}
