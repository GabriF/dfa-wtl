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

func TestComputeNext(t *testing.T) {
	m := &TwoWayDfaWtl{
		initialState: 0,
		delta: [][]int{
			{1, 2, 6, 7},
			{7, 7, 3, 7},
			{7, 7, 4, 7},
			{5, 7, 7, 7},
			{7, 5, 7, 7},
			{7, 7, 7, 0},
		},
		tau: [][]bool{
			{false, false},
			{true, true},
			{true, true},
			{false, false},
			{false, false},
			{true, true},
		},
		leftMarker:          ']',
		rightMarker:         '[',
		alphabetEnumeration: map[rune]int{'a': 0, 'b': 1, '[': 2, ']': 3},
		stateName:           []string{"q0", "qa", "qb", "pa", "pb", "pr", "accept", "reject"},
		qrCardinality:       3,
	}
	id := NewComputation(m, "abba")
	ComputeNext(m, id)

	if id.state != 1 || id.stateStr != "qa" || id.Halt {
		t.FailNow()
	}
}

func TestComputeAccept(t *testing.T) {
	m := &TwoWayDfaWtl{
		initialState: 0,
		delta: [][]int{
			{1, 2, 6, 7},
			{7, 7, 3, 7},
			{7, 7, 4, 7},
			{5, 7, 7, 7},
			{7, 5, 7, 7},
			{7, 7, 7, 0},
		},
		tau: [][]bool{
			{false, false},
			{true, true},
			{true, true},
			{false, false},
			{false, false},
			{true, true},
		},
		leftMarker:          ']',
		rightMarker:         '[',
		alphabetEnumeration: map[rune]int{'a': 0, 'b': 1, '[': 2, ']': 3},
		stateName:           []string{"q0", "qa", "qb", "pa", "pb", "pr", "accept", "reject"},
		qrCardinality:       3,
	}
	id := NewComputation(m, "abba")

	for !id.Halt {
		ComputeNext(m, id)
	}

	if (!id.Halt) || id.stateStr != "accept" {
		t.FailNow()
	}
}

func TestComputeReject(t *testing.T) {
	m := &TwoWayDfaWtl{
		initialState: 0,
		delta: [][]int{
			{1, 2, 6, 7},
			{7, 7, 3, 7},
			{7, 7, 4, 7},
			{5, 7, 7, 7},
			{7, 5, 7, 7},
			{7, 7, 7, 0},
		},
		tau: [][]bool{
			{false, false},
			{true, true},
			{true, true},
			{false, false},
			{false, false},
			{true, true},
		},
		leftMarker:          ']',
		rightMarker:         '[',
		alphabetEnumeration: map[rune]int{'a': 0, 'b': 1, '[': 2, ']': 3},
		stateName:           []string{"q0", "qa", "qb", "pa", "pb", "pr", "accept", "reject"},
		qrCardinality:       3,
	}
	id := NewComputation(m, "abbba")

	for !id.Halt {
		ComputeNext(m, id)
	}

	if (!id.Halt) || id.stateStr != "reject" {
		t.FailNow()
	}
}

func TestFromSymbolicDefinition(t *testing.T) {
	m := FromSymbolDescription(
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
	)

	t.Log(m)

	id := NewComputation(m, "aababbcc")
	for !id.Halt {
		t.Log(id.String())
		ComputeNext(m, id)
	}
	if id.stateStr != "accept" {
		t.Fatalf("Expected accept")
	}

	id = NewComputation(m, "aababbc")
	for !id.Halt {
		t.Log(id.String())
		ComputeNext(m, id)
	}
	if id.stateStr != "reject" {
		t.Fatalf("Expected reject")
	}
}
