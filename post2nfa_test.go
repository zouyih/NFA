package main

import (
	"encoding/json"
	"strings"
	"testing"
)

var post2nfaTests = []struct {
	given    string
	expected string
}{
	{"a", `{
		"start":0,
		"end":1,
		"state_map":{
			"0":{"id":0,"transitions":[0]},
			"1":{"id":1,"transitions":[]}
			},
		"transition_map":{
			"0":{"id":0,"from":0,"to":1,"input":"a"}
		}
	}`},
	{"ab|", `{
		"start":4,
		"end":5,
		"state_map":{
			"0":{"id":0,"transitions":[0]},
			"1":{"id":1,"transitions":[4]},
			"2":{"id":2,"transitions":[1]},
			"3":{"id":3,"transitions":[5]},
			"4":{"id":4,"transitions":[2,3]},
			"5":{"id":5,"transitions":[]}
			},
		"transition_map":{
		"0":{"id":0,"from":0,"to":1,"input":"a"},
		"1":{"id":1,"from":2,"to":3,"input":"b"},
		"2":{"id":2,"from":4,"to":0,"input":""},
		"3":{"id":3,"from":4,"to":2,"input":""},
		"4":{"id":4,"from":1,"to":5,"input":""},
		"5":{"id":5,"from":3,"to":5,"input":""}
		}
	}`},
	{"ab.", `{
		"start":0,
		"end":3,
		"state_map":{
			"0":{"id":0,"transitions":[0]},
			"1":{"id":1,"transitions":[1]},
			"3":{"id":3,"transitions":[]}
		},
		"transition_map":{
			"0":{"id":0,"from":0,"to":1,"input":"a"},
			"1":{"id":1,"from":1,"to":3,"input":"b"}
		}
	}`},
	{"a*", `{
		"start":2,
		"end":3,
		"state_map":{
			"0":{"id":0,"transitions":[0]},
			"1":{"id":1,"transitions":[3,4]},
			"2":{"id":2,"transitions":[1,2]},
			"3":{"id":3,"transitions":[]}
		},
		"transition_map":{
			"0":{"id":0,"from":0,"to":1,"input":"a"},
			"1":{"id":1,"from":2,"to":3,"input":""},
			"2":{"id":2,"from":2,"to":0,"input":""},
			"3":{"id":3,"from":1,"to":3,"input":""},
			"4":{"id":4,"from":1,"to":0,"input":""}
			}
		}`},
}

func TestPost2nfa(t *testing.T) {
	for _, test := range post2nfaTests {
		nfa := post2nfa(test.given)
		actual, _ := json.Marshal(nfa)
		expected := strings.Replace(test.expected, "\n", "", -1)
		expected = strings.Replace(expected, "	", "", -1)
		if string(actual) != expected {
			t.Fatalf("give %s\n expected %s\n actual %s\n", test.given, test.expected, actual)
		}
	}
}
