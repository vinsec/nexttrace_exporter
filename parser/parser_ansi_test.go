package parser

import (
	"testing"
)

func TestCleanNextTraceOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "output with ANSI codes",
			input:    "\x1b[37;1m[NextTrace API]\x1b[0;22m preferred API IP - \x1b[32;1m[2606:4700:20::681a:c97]\x1b[0;22m\n{\"Hops\":[]}",
			expected: "{\"Hops\":[]}",
		},
		{
			name:     "pure JSON without ANSI",
			input:    "{\"Hops\":[]}",
			expected: "{\"Hops\":[]}",
		},
		{
			name:     "complex ANSI with actual JSON",
			input:    "\x1b[37;1m[NextTrace API]\x1b[0;22m test\n{\"Hops\":[[{\"Success\":true}]]}",
			expected: "{\"Hops\":[[{\"Success\":true}]]}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanNextTraceOutput([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("cleanNextTraceOutput() = %q, want %q", string(result), tt.expected)
			}
		})
	}
}
