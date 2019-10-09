package main

import "testing"

type strTest struct {
	str      string
	expected bool
}

var matchTests = []struct {
	patern string
	strs   []strTest
}{
	{
		"a",
		[]strTest{
			{"a", true},
			{"b", false},
			{"", false},
			{"aa", false},
		},
	},
	{
		"ab",
		[]strTest{
			{"ab", true},
			{"a", false},
			{"b", false},
			{"ac", false},
			{"abc", false},
		},
	},
	{
		"a|b",
		[]strTest{
			{"a", true},
			{"b", true},
			{"c", false},
			{"ab", false},
			{"", false},
		},
	},
	{
		"(a|b)c",
		[]strTest{
			{"ac", true},
			{"bc", true},
			{"abc", false},
			{"c", false},
		},
	},
	{
		"a+b|c",
		[]strTest{
			{"ab", true},
			{"ac", false},
			{"b", false},
			{"aab", true},
			{"aaaaaab", true},
			{"c", true},
		},
	},
	{
		"a?bc",
		[]strTest{
			{"abc", true},
			{"bc", true},
			{"ab", false},
			{"ac", false},
		},
	},
	{
		"ab?(cd)+",
		[]strTest{
			{"abcd", true},
			{"acd", true},
			{"abcdcdcd", true},
			{"abcdcd", true},
			{"ac", false},
			{"ab", false},
		},
	},
	{
		"中文?(测试)+",
		[]strTest{
			{"中文测试", true},
			{"中测试", true},
			{"中文测试测试测试", true},
			{"中测试测试", true},
			{"中测", false},
			{"中文", false},
		},
	},
}

func TestMatch(t *testing.T) {
	for _, pattenTest := range matchTests {
		pattern := pattenTest.patern
		post := re2post(pattern)
		nfa := post2nfa(post)
		for _, strItem := range pattenTest.strs {
			actual := nfa.match(strItem.str)
			if actual != strItem.expected {
				t.Fatalf("pattern: %s, str: %s, expected: %v, actual: %v", pattern, strItem.str, strItem.expected, actual)
			}
		}
	}
}
