package engine

import (
	"testing"
)

func TestTargetSize(t *testing.T) {
	cases := []struct {
		level    int
		expected int
	}{
		{level: 0, expected: 0},
		{level: 1, expected: 6},
		{level: 5, expected: 6},
		{level: 6, expected: 8},
		{level: 7, expected: 8},
		{level: 8, expected: 9},
		{level: 9, expected: 10},
		{level: 10, expected: 10},
	}

	for _, c := range cases {
		got := TargetSize(c.level)
		if got != c.expected {
			t.Errorf("TargetSize(%d) = %d; want %d", c.level, got, c.expected)
		}
	}
}
