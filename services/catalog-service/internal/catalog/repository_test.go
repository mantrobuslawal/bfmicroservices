package catalog

import "testing"

func TestNormaliseLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{name: "default when zero", input: 0, expected: 20},
		{name: "default when negative", input: -1, expected: 20},
		{name: "keeps valid value", input: 50, expected: 50},
		{name: "caps high value", input: 500, expected: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := normaliseLimit(tt.input)
			if actual != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, actual)
			}
		})
	}
}

func TestNormaliseOffset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{name: "zero", input: 0, expected: 0},
		{name: "negative", input: -10, expected: 0},
		{name: "positive", input: 25, expected: 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := normaliseOffset(tt.input)
			if actual != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, actual)
			}
		})
	}
}
