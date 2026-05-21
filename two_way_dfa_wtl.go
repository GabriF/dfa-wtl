package gotwodfawtl

import "fmt"

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
	Halt     bool
	stateStr string
	state    int
	tape     runeDoublyLinkedList
}

func (t *TwoWayDfaWtlID) String() string {
	if t.Halt {
		return t.stateStr
	} else {
		return fmt.Sprintf("(%s, %s)", t.stateStr, t.tape.String())
	}
}

func newID(stateStr string, state int, tape runeDoublyLinkedList) *TwoWayDfaWtlID {
	isHalted := stateStr == ACCEPT_STATE || stateStr == REJECT_STATE
	return &TwoWayDfaWtlID{
		isHalted,
		stateStr,
		state,
		tape,
	}
}

func NewComputation(m *TwoWayDfaWtl, word string) *TwoWayDfaWtlID {
	return newID(m.stateName[m.initialState], m.initialState, *fromStringDLL("]" + word + "["))
}

func ComputeNext(m *TwoWayDfaWtl, id *TwoWayDfaWtlID) {
	if id.Halt {
		return
	}

	var currentSymbol *runeNode
	var nextSymbol func(*runeNode) *runeNode
	isCurrentStateQr := id.state < m.qrCardinality
	if isCurrentStateQr {
		currentSymbol = id.tape.head.next
		nextSymbol = func(rn *runeNode) *runeNode {
			return rn.next
		}
	} else {
		currentSymbol = id.tape.tail.prev
		nextSymbol = func(rn *runeNode) *runeNode {
			return rn.prev
		}
	}

	for currentSymbol.val != m.leftMarker &&
		currentSymbol.val != m.rightMarker &&
		m.tau[id.state][m.alphabetEnumeration[currentSymbol.val]] {
		currentSymbol = nextSymbol(currentSymbol)
	}

	if currentSymbol.val != m.leftMarker && currentSymbol.val != m.rightMarker {
		id.tape.remove(currentSymbol)
	}

	nextState := m.delta[id.state][m.alphabetEnumeration[currentSymbol.val]]
	*id = *newID(m.stateName[nextState], nextState, id.tape)
}
