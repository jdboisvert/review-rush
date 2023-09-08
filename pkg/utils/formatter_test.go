package utils

import (
	"testing"
)

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name     string
		counts   map[string]int
		expected string
	}{
		{
			name:     "single user",
			counts:   map[string]int{"Alice": 5},
			expected: "Code review rankings for today:\nAlice: 5 reviews\n",
		},
		{
			name:     "multiple users",
			counts:   map[string]int{"Alice": 3, "Bob": 5},
			expected: "Code review rankings for today:\nAlice: 3 reviews\nBob: 5 reviews\n",
		},
		{
			name:     "no users",
			counts:   map[string]int{},
			expected: "Code review rankings for today:\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatMessage(tt.counts)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
