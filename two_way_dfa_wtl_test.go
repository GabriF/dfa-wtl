package gotwodfawtl

import "testing"

func TestLinkedListOneElement(t *testing.T) {
	ll := fromStringDLL("a")
	if ll.String() != "a" {
		t.FailNow()
	}
}

func TestLinkedListRemove(t *testing.T) {
	ll := fromStringDLL("aba")
	toRemove := ll.head.next
	ll.remove(toRemove)
	if ll.String() != "aa" {
		t.FailNow()
	}
}

func TestLinkedListRemoveHead(t *testing.T) {
	ll := fromStringDLL("aba")
	toRemove := ll.head
	ll.remove(toRemove)
	if ll.String() != "ba" {
		t.FailNow()
	}
}

func TestLinkedListRemoveTail(t *testing.T) {
	ll := fromStringDLL("aba")
	toRemove := ll.tail
	ll.remove(toRemove)
	if ll.String() != "ab" {
		t.FailNow()
	}
}

func automatons() []*TwoWayDfaWtl {
	return []*TwoWayDfaWtl{FromSymbolDescription(
		[]string{"qa", "qb"},
		[]string{"qc"},
		"qa",
		[]rune{'a', 'b', 'c'},
		']',
		'[',
		map[string]map[rune]string{
			"qa": {
				'a': "qb",
			},
			"qb": {
				'b': "qc",
			},
			"qc": {
				'c': "qa",
				']': "accept",
			},
		},
		map[string]map[rune]bool{
			"qa": {
				'b': true,
			},
			"qb": {
				'a': true,
			},
		},
	),
		FromSymbolDescription(
			[]string{"q0", "qa", "qb"},
			[]string{"pa", "pb", "pr"},
			"q0",
			[]rune{'a', 'b'},
			']',
			'[',
			map[string]map[rune]string{
				"q0": {
					'[': "accept",
					'a': "qa",
					'b': "qb",
				},
				"qa": {
					'[': "pa",
				},
				"qb": {
					'[': "pb",
				},
				"pa": {
					'a': "pr",
				},
				"pb": {
					'b': "pr",
				},
				"pr": {
					']': "q0",
				},
			},
			map[string]map[rune]bool{
				"qa": {
					'a': true,
					'b': true,
				},
				"qb": {
					'a': true,
					'b': true,
				},
				"pr": {
					'a': true,
					'b': true,
				},
			},
		),
	}
}

func TestAutomatonComputation(t *testing.T) {
	tests := []struct {
		m      *TwoWayDfaWtl
		word   string
		accept bool
	}{
		{automatons()[0], "abaabbcc", true},
		{automatons()[0], "ababbccc", false},

		{automatons()[1], "ababbbbaba", true},
		{automatons()[1], "ababbbaba", false},
	}

	for _, tt := range tests {
		id := NewComputation(tt.m, tt.word)
		for !id.Halt {
			ComputeNext(tt.m, id)
		}
		correct := (id.stateStr == "accept" && tt.accept) || (id.stateStr == "reject" && (!tt.accept))
		if !correct {
			t.Errorf("Expected %t on input %s with automaton %+v", tt.accept, tt.word, tt.m)
		}
	}
}
