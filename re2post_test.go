package main

import "testing"

var tests = []struct {
	given    string
	expected string
}{
	{"ab", "ab."},
	{"a|b", "ab|"},
	{"a|b|c", "ab|c|"},
	{"a*b", "a*b."},
	{"a?b", "a?b."},
	{"a+b", "a+b."},
	{"ab+|cd", "ab+.cd.|"},
	{"(a|b)c", "ab|c."},
	{"abc", "ab.c."},
	{"(a|b)(c|d)", "ab|cd|."},
	{"ab?", "ab?."},
	{"(a|b)c?", "ab|c?."},
	{"(a+|b)c", "a+b|c."},
	{"a(bc)+", "abc.+."},
}

func TestRe2post(t *testing.T) {
	for _, test := range tests {
		actual := re2post(test.given)
		if actual != test.expected {
			t.Fatalf("give %s, expected %s, actual %s", test.given, test.expected, actual)
		}
	}
}
