package gotwodfawtl

import (
	"fmt"
	"strings"
)

const (
	ACCEPT_STATE = "accept"
	REJECT_STATE = "reject"
)

type TwoWayDfaWtl struct {
	initialState int

	delta [][]int
	tau   [][]bool

	leftMarker  rune
	rightMarker rune

	alphabetEnumeration map[rune]int

	stateName     []string
	qrCardinality int
}

func FromSymbolDescription(
	qrStates []string, qlStates []string, initialState string,
	sigma []rune, leftMarker rune, rightMarker rune,
	delta map[string]map[rune]string,
	tau map[string]map[rune]bool) *TwoWayDfaWtl {
	m := &TwoWayDfaWtl{}

	m.leftMarker = leftMarker
	m.rightMarker = rightMarker

	m.qrCardinality = len(qrStates)

	statesCardinality := len(qrStates) + len(qlStates)
	supportStatesEnumeration := make(map[string]int, statesCardinality+2)
	m.stateName = make([]string, statesCardinality+2)
	for i, q := range qrStates {
		m.stateName[i] = q
		supportStatesEnumeration[q] = i
	}
	for i, q := range qlStates {
		m.stateName[i+len(qrStates)] = q
		supportStatesEnumeration[q] = i + len(qrStates)
	}
	m.stateName[statesCardinality] = ACCEPT_STATE
	m.stateName[statesCardinality+1] = REJECT_STATE
	supportStatesEnumeration[ACCEPT_STATE] = statesCardinality
	supportStatesEnumeration[REJECT_STATE] = statesCardinality + 1
	m.alphabetEnumeration = make(map[rune]int, len(sigma)+2)
	for i, a := range sigma {
		m.alphabetEnumeration[a] = i
	}
	m.alphabetEnumeration[leftMarker] = len(sigma)
	m.alphabetEnumeration[rightMarker] = len(sigma) + 1

	m.initialState = supportStatesEnumeration[initialState]

	m.delta = make([][]int, statesCardinality)
	m.tau = make([][]bool, statesCardinality)
	for i := range statesCardinality {
		m.delta[i] = make([]int, len(sigma)+2)
		m.tau[i] = make([]bool, len(sigma))
		q := m.stateName[i]

		for a, j := range m.alphabetEnumeration {
			if val, ok := delta[q][a]; !ok {
				m.delta[i][j] = supportStatesEnumeration[REJECT_STATE]
			} else {
				m.delta[i][j] = supportStatesEnumeration[val]
			}
		}

		for j, a := range sigma {
			if val, ok := tau[q][a]; !ok {
				m.tau[i][j] = false
			} else {
				m.tau[i][j] = val
			}
		}
	}

	return m
}

type TwoWayDfaWtlID struct {
	Accept    bool
	Halt      bool
	stateStr  string
	state     int
	head      []int
	storeHead []int
	tail      []int
	storeTail []int
	automaton *TwoWayDfaWtl
}

func (t *TwoWayDfaWtlID) String() string {
	if t.Halt {
		return t.stateStr
	}

	tapeContent := make([]rune, len(t.storeHead))
	for i := range len(tapeContent) {
		tapeContent[i] = -1
	}

	for letter, letterEnum := range t.automaton.alphabetEnumeration {
		if letter == t.automaton.leftMarker || letter == t.automaton.rightMarker {
			continue
		}
		currentHead := t.head[letterEnum]
		for currentHead != -1 {
			tapeContent[currentHead] = letter
			currentHead = t.storeHead[currentHead]
		}
	}

	str := &strings.Builder{}
	str.WriteString(fmt.Sprintf("(%s, ", t.stateStr))
	for _, a := range tapeContent {
		if a != -1 {
			str.WriteRune(a)
		}
	}
	str.WriteRune(')')

	return str.String()
}

func NewComputation(m *TwoWayDfaWtl, word string) *TwoWayDfaWtlID {
	head := make([]int, len(m.alphabetEnumeration)-2)
	storeHead := make([]int, len(word))
	tail := make([]int, len(m.alphabetEnumeration)-2)
	storeTail := make([]int, len(word))

	for i := range len(head) {
		head[i] = -1
		tail[i] = -1
	}

	runeWord := []rune(word)
	for i := len(word) - 1; i >= 0; i-- {
		a := runeWord[i]
		storeHead[i] = head[m.alphabetEnumeration[a]]
		head[m.alphabetEnumeration[a]] = i
	}

	for i := 0; i < len(word); i++ {
		a := runeWord[i]
		storeTail[i] = tail[m.alphabetEnumeration[a]]
		tail[m.alphabetEnumeration[a]] = i
	}

	return &TwoWayDfaWtlID{
		Accept:    false,
		Halt:      false,
		stateStr:  m.stateName[m.initialState],
		state:     m.initialState,
		head:      head,
		storeHead: storeHead,
		tail:      tail,
		storeTail: storeTail,
		automaton: m,
	}
}

func (t *TwoWayDfaWtlID) changeState(letterEnum int) {
	nextState := t.automaton.delta[t.state][letterEnum]
	t.state = nextState
	t.stateStr = t.automaton.stateName[t.state]
	if t.stateStr == ACCEPT_STATE || t.stateStr == REJECT_STATE {
		t.Halt = true
		t.Accept = t.stateStr == ACCEPT_STATE
	}
}

func ComputeNext(id *TwoWayDfaWtlID) {
	if id.Halt {
		return
	}

	m := id.automaton

	var toReadLetterEnum int
	if id.state < m.qrCardinality {
		toReadLetterEnum = m.alphabetEnumeration[m.rightMarker]

		min := len(id.storeHead)
		for i, v := range id.head {
			if v != -1 && (!m.tau[id.state][i]) && v < min {
				min = v
				toReadLetterEnum = i
			}
		}
	} else {
		toReadLetterEnum = m.alphabetEnumeration[m.leftMarker]

		maxPos := 0
		for i, v := range id.tail {
			if (!m.tau[id.state][i]) && v > maxPos {
				maxPos = v
				toReadLetterEnum = i
			}
		}
	}

	if id.state < m.qrCardinality &&
		toReadLetterEnum != m.alphabetEnumeration[m.rightMarker] {
		id.head[toReadLetterEnum] = id.storeHead[id.head[toReadLetterEnum]]
		if id.head[toReadLetterEnum] == -1 {
			id.storeTail[id.tail[toReadLetterEnum]] = -1
			id.tail[toReadLetterEnum] = -1
		} else {
			id.storeTail[id.head[toReadLetterEnum]] = -1
		}
	} else if id.state >= m.qrCardinality &&
		toReadLetterEnum != m.alphabetEnumeration[m.leftMarker] {
		id.tail[toReadLetterEnum] = id.storeTail[id.tail[toReadLetterEnum]]
		if id.tail[toReadLetterEnum] == -1 {
			id.storeHead[id.head[toReadLetterEnum]] = -1
			id.head[toReadLetterEnum] = -1
		} else {
			id.storeHead[id.tail[toReadLetterEnum]] = -1
		}
	}

	id.changeState(toReadLetterEnum)
}
