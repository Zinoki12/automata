package main

import (
	"fmt"
	"unicode"
)

type class rune

const (
	ClassDigit    class = 'd'
	ClassNotDigit class = 'n'
	ClassSpace    class = 's'
	ClassComma    class = 'c'
	ClassOther    class = 'o'
)

type Key struct {
	State int
	Class class
}

type DFA struct {
	StartState  int
	FinalStates map[int]struct{}
	Transitions map[Key]int
	Classifier  func(rune) class
}

func ClassifyNonDigits(r rune) class {
	if unicode.IsDigit(r) {
		return ClassDigit
	}
	return ClassNotDigit
}

func NewNonDigitDFA() DFA {
	return DFA{
		StartState: 0,
		FinalStates: map[int]struct{}{
			1: {},
		},
		Transitions: map[Key]int{
			{0, ClassNotDigit}: 1,
			{1, ClassNotDigit}: 1,
		},
		Classifier: ClassifyNonDigits,
	}
}

func NewDigitDFA() DFA {
	return DFA{
		StartState: 0,
		FinalStates: map[int]struct{}{
			1: {},
		},
		Transitions: map[Key]int{
			{0, ClassDigit}: 1,
			{1, ClassDigit}: 1,
		},
		Classifier: ClassifyNonDigits,
	}
}

func (dfa DFA) AcceptsString(s string) bool {
	currentState := dfa.StartState
	for _, r := range s {
		currentClass := dfa.Classifier(r)
		state, ok := dfa.Transitions[Key{currentState, currentClass}]
		if !ok {
			return false
		}
		currentState = state
	}

	_, ok := dfa.FinalStates[currentState]
	return ok
}

func (dfa DFA) MatchAt(text []rune, pos int) (length int, match bool) {
	currentState := dfa.StartState
	lastMatchLen := -1

	if _, ok := dfa.FinalStates[currentState]; ok {
		lastMatchLen = 0
	}

	for i := pos; i < len(text); i++ {
		currentClass := dfa.Classifier(text[i])
		state, ok := dfa.Transitions[Key{currentState, currentClass}]
		if !ok {
			break
		}
		currentState = state
		if _, ok = dfa.FinalStates[currentState]; ok {
			lastMatchLen = i - pos + 1
		}
	}
	if lastMatchLen != -1 {
		return lastMatchLen, true
	}

	return 0, false
}

func (dfa DFA) Split(text string) []string {
	runeText := []rune(text)
	cut := 0
	res := make([]string, 0, 8)
	for i := 0; i < len(runeText); {
		if ln, ok := dfa.MatchAt(runeText, i); ok {
			res = append(res, string(runeText[cut:i]))
			i += ln
			cut = i
		} else {
			i++
		}

	}
	res = append(res, string(runeText[cut:]))
	return res
}

func run(name string, dfa DFA, inputs []string) {
	fmt.Println(name)
	for _, s := range inputs {
		fmt.Printf("  %-8q -> %t\n", s, dfa.AcceptsString(s))
	}
}

func main() {
	inputs := []string{"abc", " , ", "", "1", "123", "abc1", "a1c", "1a"}

	nND := NewNonDigitDFA()
	nD := NewDigitDFA()

	run(`Паттерн 1: \D+`, nND, inputs)
	run(`Паттерн 2: \d+`, nD, inputs)

	for i := range 4 {
		ln, bl := nND.MatchAt([]rune("12ab3"), i)

		fmt.Printf("pos: %d; len: %d; bool: %t\n", i, ln, bl)
	}
	ninputs := []string{"12ab3", "a12b", "ab", "1 One, 2 Two, 3 Three is good."}
	for _, i := range ninputs {
		fmt.Printf("%s - %q\n", i, nND.Split(i))
	}
}
