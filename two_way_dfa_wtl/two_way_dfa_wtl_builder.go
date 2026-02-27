package two_way_dfa_wtl

type TwoWayDFAwtlBuilder struct {
	stateNameToEnum map[string]state
	a               *TwoWayDFAwtl
}

func NewBuilder(rightState []string, leftState []string, startState string, leftMarker rune, rightMarker rune) *TwoWayDFAwtlBuilder {
	a := &TwoWayDFAwtl{}
	a.direction = map[state]stateHeadDirection{}
	a.delta = map[state]map[rune]state{}
	a.translucent = map[state]map[rune]bool{}
	a.stateName = map[state]string{}

	builder := &TwoWayDFAwtlBuilder{
		map[string]state{},
		a,
	}
	builder.stateNameToEnum["accept"] = -1

	for i, name := range rightState {
		stateEnum := state(i)
		a.stateName[stateEnum] = name
		a.direction[stateEnum] = R
		builder.stateNameToEnum[name] = stateEnum
	}

	for i, name := range leftState {
		stateEnum := state(i + len(rightState))
		a.stateName[stateEnum] = name
		a.direction[stateEnum] = L
		builder.stateNameToEnum[name] = stateEnum
	}

	a.leftMarker = leftMarker
	a.rightMarker = rightMarker
	a.start = builder.stateNameToEnum[startState]

	return builder
}

func (b *TwoWayDFAwtlBuilder) WithDelta(fromState string, letter rune, destinationState string) {
	stateEnum := b.stateNameToEnum[fromState]
	destinationStateEnum := b.stateNameToEnum[destinationState]

	if b.a.delta[stateEnum] == nil {
		b.a.delta[stateEnum] = map[rune]state{}
	}

	b.a.delta[stateEnum][letter] = destinationStateEnum
}

func (b *TwoWayDFAwtlBuilder) WithTau(fromState string, translucents ...rune) {
	stateEnum := b.stateNameToEnum[fromState]

	if b.a.translucent[stateEnum] == nil {
		b.a.translucent[stateEnum] = map[rune]bool{}
	}

	for _, letter := range translucents {
		b.a.translucent[stateEnum][letter] = true
	}
}

func (b *TwoWayDFAwtlBuilder) Build() (*TwoWayDFAwtl, bool) {
	returned := b.a
	b.a = nil
	b.stateNameToEnum = nil
	return returned, returned != nil
}
