package parrot

import (
	"strings"
	"testing"
)

type split_testpair struct {
	message  string
	expected [][]rune
}

var split_tests = []split_testpair{
	{
		"",
		[][]rune{[]rune("")},
	},
	{
		"A simple text message.",
		[][]rune{
			[]rune("A simple text message."),
		},
	},
	{
		strings.Repeat("1234567890", 16),
		[][]rune{
			[]rune(strings.Repeat("1234567890", 16)),
		},
	},
	{
		strings.Repeat("1234567890", 16) + "X",
		[][]rune{
			[]rune(strings.Repeat("1234567890", 16)),
			[]rune("X"),
		},
	},
	// todo test with unicodes
}

func checkSlicesEqual(a, b []rune) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func checkSlicesOfSlicesEqual(a, b [][]rune) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !checkSlicesEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestSplitMessageIntoParts(t *testing.T) {
	for _, pair := range split_tests {
		res := SplitMessageIntoParts(pair.message)
		if !checkSlicesOfSlicesEqual(res, pair.expected) {
			t.Errorf("Got: %s, Expected: %s", res, pair.expected)
		}
	}
}
