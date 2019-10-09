package main

import (
	"encoding/json"
	"fmt"
)

type State struct {
	Id          int   `json:"id"`
	Transitions []int `json:"transitions"`
}

type Transition struct {
	Id    int    `json:"id"`
	From  int    `json:"from"`
	To    int    `json:"to"`
	Input string `json:"input"`
}

type Frag struct {
	start *State
	end   *State
}

type NFA struct {
	Start         int                 `json:"start"`
	End           int                 `json:"end"`
	StateMap      map[int]*State      `json:"state_map"`
	TransitionMap map[int]*Transition `json:"transition_map"`

	numState      int
	numTransition int
}

func newNFA() NFA {
	nfa := NFA{}
	nfa.StateMap = make(map[int]*State)
	nfa.TransitionMap = make(map[int]*Transition)

	return nfa
}

func (nfa *NFA) newState() *State {
	s := new(State)
	s.Transitions = make([]int, 0)
	id := nfa.numState
	s.Id = id
	nfa.StateMap[id] = s

	nfa.numState++
	return s
}

func (nfa *NFA) newTransition(fromState, toState *State, input string) *Transition {
	id := nfa.numTransition
	transition := &Transition{id, fromState.Id, toState.Id, input}
	nfa.TransitionMap[id] = transition
	fromState.Transitions = append(fromState.Transitions, id)

	nfa.numTransition++
	return transition
}

func newFrag(start, end *State) Frag {
	return Frag{start, end}
}

func post2nfa(postfix string) NFA {
	nfa := newNFA()

	fragStack := make([]Frag, 0)

	pop := func() Frag {
		fr := fragStack[len(fragStack)-1]
		fragStack = fragStack[:len(fragStack)-1]
		return fr
	}
	push := func(fr Frag) {
		fragStack = append(fragStack, fr)
	}

	for _, char := range postfix {
		switch char {
		default:
			start := nfa.newState()
			end := nfa.newState()

			nfa.newTransition(start, end, string(char))

			fr := newFrag(start, end)
			push(fr)
		case '.':
			f2 := pop()
			f1 := pop()
			for _, transId := range f2.start.Transitions {
				transition := nfa.TransitionMap[transId]
				transition.From = f1.end.Id
				f1.end.Transitions = append(f1.end.Transitions, transId)
			}
			delete(nfa.StateMap, f2.start.Id)
			fr := newFrag(f1.start, f2.end)
			push(fr)
		case '|':
			f2 := pop()
			f1 := pop()
			start := nfa.newState()
			end := nfa.newState()

			nfa.newTransition(start, f1.start, "")
			nfa.newTransition(start, f2.start, "")
			nfa.newTransition(f1.end, end, "")
			nfa.newTransition(f2.end, end, "")

			fr := newFrag(start, end)
			push(fr)
		case '?':
			f := pop()
			start := nfa.newState()
			end := nfa.newState()

			nfa.newTransition(start, end, "")
			nfa.newTransition(start, f.start, "")
			nfa.newTransition(f.end, end, "")

			fr := newFrag(start, end)
			push(fr)
		case '*':
			f := pop()
			start := nfa.newState()
			end := nfa.newState()

			nfa.newTransition(start, end, "")
			nfa.newTransition(start, f.start, "")
			nfa.newTransition(f.end, end, "")
			nfa.newTransition(f.end, f.start, "")

			fr := newFrag(start, end)
			push(fr)
		case '+':
			f := pop()
			start := nfa.newState()
			end := nfa.newState()

			nfa.newTransition(start, f.start, "")
			nfa.newTransition(f.end, end, "")
			nfa.newTransition(f.end, f.start, "")

			fr := newFrag(start, end)
			push(fr)
		}
	}

	f := pop()
	nfa.Start = f.start.Id
	nfa.End = f.end.Id

	return nfa
}

func (nfa *NFA) match(str string) bool {
	curStates := make(map[*State]bool)
	nfa.addState(curStates, nfa.StateMap[nfa.Start])

	for _, char := range str {
		curStates = nfa.step(char, curStates)
	}

	for state, _ := range curStates {
		if state.Id == nfa.End {
			return true
		}
	}

	return false
}

func (nfa *NFA) step(char rune, curStates map[*State]bool) map[*State]bool {
	nextStates := make(map[*State]bool)
	for state, _ := range curStates {
		if len(state.Transitions) != 1 {
			continue
		}

		transition := nfa.TransitionMap[state.Transitions[0]]
		if transition.Input == string(char) {
			nextState := nfa.StateMap[transition.To]
			nfa.addState(nextStates, nextState)
		}
	}

	return nextStates
}

func (nfa *NFA) addState(states map[*State]bool, s *State) {
	if ok, _ := states[s]; ok {
		return
	}

	transitions := s.Transitions

	//end state: len(transitions) = 0
	if len(transitions) == 0 || nfa.TransitionMap[transitions[0]].Input != "" {
		states[s] = true
		return
	}

	for _, transId := range s.Transitions {
		nextStateId := nfa.TransitionMap[transId].To
		nextState := nfa.StateMap[nextStateId]
		nfa.addState(states, nextState)
	}

}

func main() {
	s := "(ab)+"
	fmt.Println(s)
	post := re2post(s)
	fmt.Println(post)
	nfa := post2nfa(post)

	b, _ := json.Marshal(nfa)
	fmt.Printf("%s\n", b)

	fmt.Println(nfa.match("ababab"))
}
