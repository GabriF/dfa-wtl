package gotwodfawtl

import "fmt"

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
	isHalted := stateStr == "accept" || stateStr == "reject"
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

	for m.tau[id.state][m.alphabetEnumeration[currentSymbol.val]] &&
		currentSymbol.val != m.leftMarker &&
		currentSymbol.val != m.rightMarker {
		currentSymbol = nextSymbol(currentSymbol)
	}

	if currentSymbol.val != m.leftMarker && currentSymbol.val != m.rightMarker {
		id.tape.remove(currentSymbol)
	}

	nextState := m.delta[id.state][m.alphabetEnumeration[currentSymbol.val]]
	*id = *newID(m.stateName[nextState], nextState, id.tape)
}
