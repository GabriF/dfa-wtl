package gotwodfawtl

import "fmt"

const (
	_REJECT = -1
	_ACCEPT = -2
)

type TwoWayDfaWtl struct {
	initialState int

	leftMarker  rune
	rightMarker rune

	delta [][]int
	tau   [][]bool

	alphabetEnumeration map[rune]int

	stateName     []string
	qrCardinality int
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
	return &TwoWayDfaWtlID{
		false,
		stateStr,
		state,
		tape,
	}
}

func NewComputation(m *TwoWayDfaWtl, word string) *TwoWayDfaWtlID {
	return newID(m.stateName[m.initialState], m.initialState, *fromStringDLL(word))
}

func ComputeNext(m *TwoWayDfaWtl, id *TwoWayDfaWtlID) {
	if id.Halt {
		return
	}

	var currentSymbol *runeNode
	var nextSymbol func(*runeNode) *runeNode
	isCurrentStateQr := id.state < m.qrCardinality
	if isCurrentStateQr {
		currentSymbol = id.tape.head
		nextSymbol = func(rn *runeNode) *runeNode {
			return rn.next
		}
	} else {
		currentSymbol = id.tape.tail
		nextSymbol = func(rn *runeNode) *runeNode {
			return rn.prev
		}
	}

	for currentSymbol != nil &&
		m.tau[id.state][m.alphabetEnumeration[currentSymbol.val]] {
		currentSymbol = nextSymbol(currentSymbol)
	}

	var nextState int
	if currentSymbol == nil {
		if isCurrentStateQr {
			nextState = m.delta[id.state][m.alphabetEnumeration[m.rightMarker]]
		} else {
			nextState = m.delta[id.state][m.alphabetEnumeration[m.leftMarker]]
		}
	} else {
		id.tape.remove(currentSymbol)
		nextState = m.delta[id.state][m.alphabetEnumeration[currentSymbol.val]]
	}

	switch nextState {
	case _ACCEPT:
		id.Halt = true
		id.stateStr = "accept"
	case _REJECT:
		id.Halt = true
		id.stateStr = "reject"
	default:
		id.state = nextState
		id.stateStr = m.stateName[nextState]
	}
}
