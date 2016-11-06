package parrot

import (
	"reflect"
	"strings"
	"testing"
)

type split_testargs struct {
	message  string
	expected []string
}

var split_tests = []split_testargs{
	{
		"",
		[]string{""},
	},
	{
		"A simple text message.",
		[]string{"A simple text message."},
	},
	{
		strings.Repeat("1234567890", 16),
		[]string{strings.Repeat("1234567890", 16)},
	},
	{
		strings.Repeat("1234567890", 16) + "X",
		[]string{strings.Repeat("1234567890", 16), "X"},
	},
}

func TestSplitMessageIntoParts(t *testing.T) {
	for _, args := range split_tests {
		res, err := splitMessageIntoParts(args.message, 160)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(res, args.expected) {
			t.Errorf("got: %s, expected: %s", res, args.expected)
		}
	}
}

type udhheader_testargs struct {
	index     int
	total     int
	ref       byte
	expected  string
	expextErr bool
}

var udhheader_tests = []udhheader_testargs{
	{0, 1, byte(42), "0500032a0101", false},
	{1, 5, byte(42), "0500032a0502", false},
	{256, 5, byte(42), "", true},
	{3, 256, byte(42), "", true},
	{42, 42, byte(42), "", true},
}

func TestGenerateUDHHeader(t *testing.T) {
	for _, args := range udhheader_tests {
		res, err := generateUDHHeader(args.index, args.total, args.ref)
		if args.expextErr {
			if err == nil {
				t.Error("error not raised")
			}
		} else {
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(res, args.expected) {
				t.Errorf("Got: %s, Expected: %s", res, args.expected)
			}
		}

	}
}
